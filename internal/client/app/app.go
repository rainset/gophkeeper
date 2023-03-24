package app

import (
	"embed"
	_ "embed"
	"errors"
	"github.com/rainset/gophkeeper/pkg/crypt"
	"image/color"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/internal/client/model"
	"github.com/rainset/gophkeeper/internal/client/service"
	"github.com/rainset/gophkeeper/internal/client/service/channel"
	"github.com/rainset/gophkeeper/internal/client/storage"
	"github.com/rainset/gophkeeper/internal/client/ui"
	smodel "github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/pkg/hash"
	"github.com/rainset/gophkeeper/pkg/logger"
)

//go:embed icons/*
var icons embed.FS

type App struct {
	window      fyne.Window
	cfg         *config.Config
	db          *storage.Base
	HTTPService *service.HTTPService
	FileService *service.FileService
	Channels    channel.Channels
}

func New(cfg *config.Config) *App {
	fyneApp := app.New()
	w := fyneApp.NewWindow("Gophkeeper")
	w.Resize(fyne.NewSize(500, 700))

	w.SetIcon(theme.FyneLogo())

	// создаем директорию для файлов пользователя
	if _, err := os.Stat(cfg.ClientFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(cfg.ClientFolder, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	// подключение бд клиента
	db, err := storage.New(filepath.Join(cfg.ClientFolder, "gophkeeper.db"))
	if err != nil {
		log.Fatal(err)
	}

	fileService, err := service.New(filepath.Join(cfg.ClientFolder, "upload"))
	if err != nil {
		log.Fatal(err)
	}

	HTTPService := service.NewHTTPService(cfg)

	return &App{
		window:      w,
		cfg:         cfg,
		db:          db,
		HTTPService: HTTPService,
		FileService: fileService,
		Channels: channel.Channels{
			SyncProgressBar:     make(chan float64, 1),
			SyncProgressBarQuit: make(chan bool, 1),
		},
	}
}

func (a *App) Run() {
	a.pageAuth()
	a.window.ShowAndRun()
}

func (a *App) pageAuth() {

	authForm := a.authForm()
	regForm := a.regForm()

	tabs := container.NewAppTabs(
		container.NewTabItem("Авторизация", container.New(layout.NewPaddedLayout(), authForm)),
		container.NewTabItem("Регистрация", container.New(layout.NewPaddedLayout(), regForm)),
	)

	logoIcon, err := icons.ReadFile("icons/logo.png")
	if err != nil {
		logger.Error(err)
	}
	logoResource := canvas.NewImageFromResource(fyne.NewStaticResource("logo", logoIcon))
	logoResource.SetMinSize(fyne.NewSize(500, 150))
	a.window.SetContent(container.NewVBox(logoResource, tabs))
}

func (a *App) pageMain(dataType ui.DataType) {
	var content *fyne.Container
	var tabs *container.AppTabs

	syncBar := widget.NewProgressBar()
	syncBar.Hide()

	go func() {
		for {
			select {
			case v := <-a.Channels.SyncProgressBar:
				syncBar.SetValue(v)
			case quit := <-a.Channels.SyncProgressBarQuit:
				if quit {
					syncBar.Hide()
					tabs.Show()
					a.pageMain(dataType)
				}
			}
		}
	}()

	addBtn := widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
		selectedTab := tabs.Selected()
		switch selectedTab.Text {
		case ui.TabCard.String():
			dataType = ui.TypeCard
		case ui.TabCred.String():
			dataType = ui.TypeCred
		case ui.TabText.String():
			dataType = ui.TypeText
		case ui.TabFile.String():
			dataType = ui.TypeFile
		}
		a.pageAdd(0, dataType)
	})

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Синхронизация", theme.StorageIcon(), func() {
			syncBar.SetValue(0)
			syncBar.Show()
			tabs.Hide()
			a.SyncData()
		}),
		addBtn,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Выйти", theme.ContentClearIcon(), func() {
			a.pageAuth()
		}),
	)

	tabCard := container.NewTabItem(ui.TabCard.String(), container.New(layout.NewPaddedLayout(), a.cardList()))
	tabCred := container.NewTabItem(ui.TabCred.String(), container.New(layout.NewPaddedLayout(), a.credList()))
	tabText := container.NewTabItem(ui.TabText.String(), container.New(layout.NewPaddedLayout(), a.textList()))
	tabFile := container.NewTabItem(ui.TabFile.String(), container.New(layout.NewPaddedLayout(), a.fileList()))

	tabs = container.NewAppTabs(
		tabCard,
		tabCred,
		tabText,
		tabFile,
	)

	switch dataType {
	case ui.TypeCard:
		tabs.Select(tabCard)
	case ui.TypeCred:
		tabs.Select(tabCred)
	case ui.TypeText:
		tabs.Select(tabText)
	case ui.TypeFile:
		tabs.Select(tabFile)
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		syncBar,
		tabs,
	)

	a.window.SetContent(content)
}

func (a *App) pageAdd(localID int, dataType ui.DataType) {
	logger.Infof("pageadd %d %s", localID, dataType)

	var content *fyne.Container

	barTitle := canvas.NewText("Новая запись", color.Black)

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
			a.pageMain(dataType)
		}),
		layout.NewSpacer(),
		barTitle,
	)

	tabCard := container.NewTabItem(ui.TabCard.String(), container.New(layout.NewPaddedLayout(), a.addCardForm(localID)))
	tabCred := container.NewTabItem(ui.TabCred.String(), container.New(layout.NewPaddedLayout(), a.addCredForm(localID)))
	tabText := container.NewTabItem(ui.TabText.String(), container.New(layout.NewPaddedLayout(), a.addTextForm(localID)))
	tabFile := container.NewTabItem(ui.TabFile.String(), container.New(layout.NewPaddedLayout(), a.addFileForm(localID)))

	tabs := container.NewAppTabs(
		tabCard,
		tabCred,
		tabText,
		tabFile,
	)

	switch dataType {
	case ui.TypeCard:
		tabs.Select(tabCard)
	case ui.TypeCred:
		tabs.Select(tabCred)
	case ui.TypeText:
		tabs.Select(tabText)
	case ui.TypeFile:
		tabs.Select(tabFile)
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		tabs,
	)

	a.window.SetContent(content)
}

func (a *App) pageEdit(localID int, dataType ui.DataType) {

	log.Println(a.GetAllFiles())
	logger.Infof("pageEdit %d %s", localID, dataType)

	var content *fyne.Container

	deleteBtn := widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Удалить", "Вы действительно хотите удалить запись?", func(b bool) {
			if b {
				var err error
				switch dataType {
				case ui.TypeCard:
					err = a.DeleteCard(localID)
				case ui.TypeCred:
					err = a.DeleteCred(localID)
				case ui.TypeText:
					err = a.DeleteText(localID)
				case ui.TypeFile:
					err = a.DeleteFile(localID)
				}

				if err != nil {

					logger.Error("delete:", err, dataType, localID)
					dialog.ShowError(errors.New("ошибка при удалении записи"), a.window)

					return
				}

				a.pageMain(dataType)
			}
		}, a.window)
	})

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
			logger.Error("Назад ", dataType)
			a.pageMain(dataType)
		}),
		layout.NewSpacer(),
		deleteBtn,
	)

	var editCont *fyne.Container

	editCont = container.NewCenter()

	switch dataType {
	case ui.TypeCard:
		editCont = container.New(layout.NewPaddedLayout(), a.addCardForm(localID))
	case ui.TypeCred:
		editCont = container.New(layout.NewPaddedLayout(), a.addCredForm(localID))
	case ui.TypeText:
		editCont = container.New(layout.NewPaddedLayout(), a.addTextForm(localID))
	case ui.TypeFile:
		editCont = container.New(layout.NewPaddedLayout(), a.addFileForm(localID))
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		editCont,
	)

	a.window.SetContent(content)
}

func (a *App) authForm() *widget.Form {
	login := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	login.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	pass.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	authForm := widget.NewForm(
		widget.NewFormItem("Логин", login),
		widget.NewFormItem("Пароль", pass),
	)

	authForm.SubmitText = "Войти"
	authForm.OnSubmit = func() {
		a.SetUser(login.Text)

		c, err := a.GetUserConfig()
		if err != nil {
			dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)

			return
		}

		tokens, err := a.HTTPService.SignIn(model.User{Login: login.Text, Password: pass.Text})
		if err != nil {
			if errors.Is(err, service.ErrStatusUnauthorized) {
				dialog.ShowError(service.ErrStatusUnauthorized, a.window)

				return
			}

			if hash.Sha256(pass.Text) != c.Password {
				dialog.ShowError(service.ErrStatusUnauthorized, a.window)

				return
			}
		}

		c.Login = login.Text
		c.Password = hash.Sha256(pass.Text)
		c.RefreshToken = tokens.RefreshToken
		c.AccessToken = tokens.AccessToken

		signKey, err := a.HTTPService.GetSignKey(tokens.AccessToken, login.Text, pass.Text)
		if err != nil {
			dialog.ShowError(errors.New("ошибка получения ключа подписи"), a.window)

			return
		}

		c.SignKey = signKey

		err = a.SetUserConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)

			return
		}

		a.pageMain(ui.TypeCard)
		//a.SyncData()
	}

	return authForm
}

func (a *App) regForm() *widget.Form {
	login := widget.NewEntry()
	pass := widget.NewPasswordEntry()
	pass2 := widget.NewPasswordEntry()

	login.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	pass.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	pass2.Validator = validation.NewAllStrings(
		validation.NewRegexp("^.{1,}", "обязательное поле"),
		func(string) error {
			if pass.Text != pass2.Text {
				return errors.New("пароли не совпадают")
			}

			return nil
		},
	)

	regForm := widget.NewForm(
		widget.NewFormItem("Логин", login),
		widget.NewFormItem("Пароль", pass),
		widget.NewFormItem("Повторите пароль", pass2),
	)

	regForm.SubmitText = "Зарегистрироваться"
	regForm.OnSubmit = func() {
		tokens, err := a.HTTPService.SignUp(model.User{Login: login.Text, Password: pass.Text})
		if err != nil {
			logger.Error(err)

			if errors.Is(err, service.ErrStatusLoginExists) {
				dialog.ShowError(service.ErrStatusLoginExists, a.window)

				return
			}

			dialog.ShowError(service.ErrServer, a.window)

			return
		}

		a.SetUser(login.Text)
		c, err := a.GetUserConfig()
		if err != nil {
			dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)

			return
		}

		c.Login = login.Text
		c.Password = hash.Sha256(pass.Text)
		c.RefreshToken = tokens.RefreshToken
		c.AccessToken = tokens.AccessToken

		signKey, err := a.HTTPService.GetSignKey(tokens.AccessToken, login.Text, pass.Text)
		if err != nil {
			dialog.ShowError(errors.New("ошибка получения ключа подписи"), a.window)

			return
		}

		c.SignKey = signKey

		err = a.SetUserConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)

			return
		}

		a.pageMain(ui.TypeCard)
		//a.SyncData()
	}

	return regForm
}

func (a *App) addCardForm(localID int) *fyne.Container {
	var err error
	var item model.DataCard

	var title *widget.Entry
	var cardNumber *widget.Entry
	var cardDate *widget.Entry
	var cardCvv *widget.Entry

	var meta *widget.Entry

	var addForm *widget.Form

	var titleFormItem *widget.FormItem
	var cardNumberFormItem *widget.FormItem
	var cardDateFormItem *widget.FormItem
	var cardCvvFormItem *widget.FormItem
	var metaFormItem *widget.FormItem

	title = widget.NewEntry()
	cardNumber = widget.NewEntry()
	cardDate = widget.NewEntry()
	cardCvv = widget.NewEntry()
	meta = widget.NewMultiLineEntry()

	if localID > 0 {
		item, err = a.GetCard(localID)
		if err != nil {
			logger.Error(err)
		}

		title.Text = item.Title
		cardNumber.Text = item.Number
		cardDate.Text = item.Date
		cardCvv.Text = item.Cvv
		meta.Text = item.Meta
	}

	title.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	cardNumber.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	cardDate.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	cardCvv.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	titleFormItem = widget.NewFormItem("Заголовок", title)
	cardNumberFormItem = widget.NewFormItem("Номер карты", cardNumber)
	cardDateFormItem = widget.NewFormItem("Месяц/Год", cardDate)
	cardCvvFormItem = widget.NewFormItem("CVV", cardCvv)
	metaFormItem = widget.NewFormItem("Дополнительно", meta)

	addForm = widget.NewForm(

		titleFormItem,
		cardNumberFormItem,
		cardDateFormItem,
		cardCvvFormItem,
		metaFormItem,
	)

	addForm.CancelText = "Отмена"
	addForm.OnCancel = func() {
		a.pageMain(ui.TypeCard)
	}
	addForm.SubmitText = "Сохранить"
	addForm.OnSubmit = func() {
		var err error

		cardData := model.DataCard{
			Title:     title.Text,
			Number:    cardNumber.Text,
			Date:      cardDate.Text,
			Cvv:       cardCvv.Text,
			Meta:      meta.Text,
			UpdatedAt: time.Now(),
		}

		if localID > 0 {
			cardData.LocalID = item.LocalID
			cardData.ExternalID = item.ExternalID
		}

		err = a.AddCard(&cardData, false)

		if err != nil {
			logger.Error(err)
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)

			return
		}

		a.pageMain(ui.TypeCard)
	}
	return container.NewVBox(addForm)
}

func (a *App) addCredForm(localID int) *fyne.Container {
	var err error
	var item model.DataCred

	var title *widget.Entry
	var username *widget.Entry
	var password *widget.Entry
	var meta *widget.Entry
	var addForm *widget.Form

	var titleFormItem *widget.FormItem
	var usernameFormItem *widget.FormItem
	var passwordFormItem *widget.FormItem
	var metaFormItem *widget.FormItem

	title = widget.NewEntry()
	username = widget.NewEntry()
	password = widget.NewPasswordEntry()
	meta = widget.NewMultiLineEntry()

	if localID > 0 {
		item, err = a.GetCred(localID)
		if err != nil {
			logger.Error(err)
		}
		title.Text = item.Title
		username.Text = item.Username
		password.Text = item.Password
		meta.Text = item.Meta
	}

	title.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	username.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	password.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	titleFormItem = widget.NewFormItem("Заголовок", title)
	usernameFormItem = widget.NewFormItem("Имя пользователя", username)
	passwordFormItem = widget.NewFormItem("Пароль", password)
	metaFormItem = widget.NewFormItem("Дополнительно", meta)

	addForm = widget.NewForm(
		titleFormItem,
		usernameFormItem,
		passwordFormItem,
		metaFormItem,
	)

	addForm.CancelText = "Отмена"
	addForm.OnCancel = func() {
		a.pageMain(ui.TypeCred)
	}
	addForm.SubmitText = "Сохранить"
	addForm.OnSubmit = func() {
		var err error

		credData := model.DataCred{
			Title:     title.Text,
			Username:  username.Text,
			Password:  password.Text,
			Meta:      meta.Text,
			UpdatedAt: time.Now(),
		}

		if localID > 0 {
			credData.LocalID = item.LocalID
			credData.ExternalID = item.ExternalID
		}

		err = a.AddCred(&credData, false)

		if err != nil {
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)

			return
		}

		a.pageMain(ui.TypeCred)
	}
	return container.NewVBox(addForm)
}

func (a *App) addTextForm(localID int) *fyne.Container {
	var err error
	var item model.DataText

	var title *widget.Entry
	var text *widget.Entry
	var meta *widget.Entry
	var addForm *widget.Form

	var titleFormItem *widget.FormItem
	var textFormItem *widget.FormItem
	var metaFormItem *widget.FormItem

	title = widget.NewEntry()
	text = widget.NewMultiLineEntry()
	meta = widget.NewMultiLineEntry()

	if localID > 0 {
		item, err = a.GetText(localID)
		if err != nil {
			logger.Error(err)
		}
		title.Text = item.Title
		text.Text = item.Text
		meta.Text = item.Meta
	}

	title.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	text.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	titleFormItem = widget.NewFormItem("Заголовок", title)
	textFormItem = widget.NewFormItem("Текст", text)
	metaFormItem = widget.NewFormItem("Дополнительно", meta)

	addForm = widget.NewForm(
		titleFormItem,
		textFormItem,
		metaFormItem,
	)

	addForm.CancelText = "Отмена"
	addForm.OnCancel = func() {
		a.pageMain(ui.TypeText)
	}
	addForm.SubmitText = "Сохранить"
	addForm.OnSubmit = func() {
		var err error

		textData := model.DataText{
			Title:     title.Text,
			Text:      text.Text,
			Meta:      meta.Text,
			UpdatedAt: time.Now(),
		}

		if localID > 0 {
			textData.LocalID = item.LocalID
			textData.ExternalID = item.ExternalID
		}

		err = a.AddText(&textData, false)

		if err != nil {
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)

			return
		}

		a.pageMain(ui.TypeText)
	}
	return container.NewVBox(addForm)
}

func (a *App) addFileForm(localID int) *fyne.Container {
	type TempFile struct {
		exists   bool
		filename string
		ext      string
		r        fyne.URIReadCloser
	}

	var err error
	var item model.DataFile

	var title *widget.Entry
	var uploadBtn *widget.Button
	var meta *widget.Entry
	var addForm *widget.Form

	var titleFormItem *widget.FormItem
	var uploadFormItem *widget.FormItem
	var metaFormItem *widget.FormItem
	var tempFile TempFile

	title = widget.NewEntry()
	meta = widget.NewMultiLineEntry()

	uploadBtn = widget.NewButton("Выбрать файл", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil {
				return
			}
			tempFile.exists = true
			tempFile.ext = r.URI().Extension()
			tempFile.filename = r.URI().Name()
			tempFile.r = r
			uploadBtn.SetText(r.URI().Name())
		}, a.window)
	})

	if localID > 0 {
		item, err = a.GetFile(localID)
		if err != nil {
			logger.Error(err)
		}

		title.Text = item.Title
		uploadBtn.Text = item.Filename
		meta.Text = item.Meta
	}

	title.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

	titleFormItem = widget.NewFormItem("Заголовок", title)
	uploadFormItem = widget.NewFormItem("Файл", uploadBtn)
	metaFormItem = widget.NewFormItem("Дополнительно", meta)

	addForm = widget.NewForm(
		titleFormItem,
		uploadFormItem,
		metaFormItem,
	)

	addForm.CancelText = "Отмена"
	addForm.OnCancel = func() {
		a.pageMain(ui.TypeFile)
	}
	addForm.SubmitText = "Сохранить"
	addForm.OnSubmit = func() {

		saveItem := model.DataFile{
			Title:     title.Text,
			Filename:  item.Filename,
			Path:      item.Path,
			Ext:       item.Ext,
			Meta:      meta.Text,
			UpdatedAt: time.Now(),
		}

		if tempFile.exists {
			filePath, err := a.FileService.SaveFile(tempFile.r, tempFile.ext)
			if err != nil {
				dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
				return
			}

			saveItem.Path = filePath
			saveItem.Ext = tempFile.r.URI().Extension()

			err = a.FileService.DeleteFile(item.Path)
			if err != nil {
				dialog.ShowError(errors.New("ошибка удаления старой версии файла"), a.window)
				return
			}
		}

		if localID > 0 {
			saveItem.LocalID = item.LocalID
			saveItem.ExternalID = item.ExternalID
		}

		err = a.AddFile(&saveItem, false)
		if err != nil {
			logger.Error("AddFile: ", err)
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
			return
		}

		a.pageMain(ui.TypeFile)
	}

	if localID > 0 {
		fileDir := filepath.Dir(item.Path)
		u, errP := url.Parse(fileDir)
		if errP != nil {
			logger.Error("url parse ", err)
		}

		link := widget.NewHyperlinkWithStyle("Просмотр файла", u, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})
		linkContainer := container.NewVBox(widget.NewLabel(""), canvas.NewLine(color.Black), link)

		//return container.NewVBox(addForm, linkContainer, container.NewCenter(widget.NewLabel(item.Path)))
		return container.NewVBox(addForm, linkContainer)
	}

	return container.NewVBox(addForm)
}

func (a *App) cardList() *fyne.Container {
	var err error
	var cards []model.DataCard
	var cardsList *widget.List

	noItems := container.NewCenter(canvas.NewText("Нет записей", color.Black))

	cards, err = a.GetAllCards()
	if err != nil {
		logger.Error(err)
	}
	cardsList = widget.NewList(
		func() int {
			return len(cards)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("Default")
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(cards[lii].Title)
		},
	)

	cardsList.OnSelected = func(id widget.ListItemID) {
		a.pageEdit(cards[id].LocalID, ui.TypeCard)
	}

	scroll := container.NewScroll(cardsList)
	scroll.SetMinSize(fyne.NewSize(100, 700))

	if len(cards) != 0 {
		noItems.Hide()
	}
	return container.New(layout.NewPaddedLayout(), scroll, noItems)
}

func (a *App) credList() *fyne.Container {
	var err error
	var creds []model.DataCred
	var credsList *widget.List

	noItems := container.NewCenter(canvas.NewText("Нет записей", color.Black))

	creds, err = a.GetAllCreds()
	if err != nil {
		logger.Error(err)
	}
	credsList = widget.NewList(
		func() int {
			return len(creds)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("Default")
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(creds[lii].Title)
		},
	)

	credsList.OnSelected = func(id widget.ListItemID) {
		a.pageEdit(creds[id].LocalID, ui.TypeCred)
	}

	scroll := container.NewScroll(credsList)
	scroll.SetMinSize(fyne.NewSize(100, 700))

	if len(creds) != 0 {
		noItems.Hide()
	}
	return container.New(layout.NewPaddedLayout(), scroll, noItems)
}

func (a *App) textList() *fyne.Container {
	var err error
	var texts []model.DataText
	var textsList *widget.List

	noItems := container.NewCenter(canvas.NewText("Нет записей", color.Black))

	texts, err = a.GetAllTexts()
	if err != nil {
		logger.Error(err)
	}
	textsList = widget.NewList(
		func() int {
			return len(texts)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("Default")
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(texts[lii].Title)
		},
	)

	textsList.OnSelected = func(id widget.ListItemID) {
		a.pageEdit(texts[id].LocalID, ui.TypeText)
	}

	scroll := container.NewScroll(textsList)
	scroll.SetMinSize(fyne.NewSize(100, 700))

	if len(texts) != 0 {
		noItems.Hide()
	}
	return container.New(layout.NewPaddedLayout(), scroll, noItems)
}

func (a *App) fileList() *fyne.Container {
	var err error
	var files []model.DataFile
	var filesList *widget.List

	noItems := container.NewCenter(canvas.NewText("Нет записей", color.Black))

	files, err = a.GetAllFiles()
	if err != nil {
		logger.Error(err)
	}
	filesList = widget.NewList(
		func() int {
			return len(files)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("Default")
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(files[lii].Title)
		},
	)

	filesList.OnSelected = func(id widget.ListItemID) {
		a.pageEdit(files[id].LocalID, ui.TypeFile)
	}

	scroll := container.NewScroll(filesList)
	scroll.SetMinSize(fyne.NewSize(100, 700))

	if len(files) != 0 {
		noItems.Hide()
	}
	return container.New(layout.NewPaddedLayout(), scroll, noItems)
}

func (a *App) SetUser(login string) {
	a.db.SetUser(login)
}

func (a *App) SetUserConfig(config model.UserConfig) (err error) {
	err = a.db.SetUserConfig(config)
	return err
}

func (a *App) GetUserConfig() (c model.UserConfig, err error) {
	c, err = a.db.GetUserConfig()
	return c, err
}

func (a *App) AddCard(card *model.DataCard, encrypted bool) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	if !encrypted {
		sKey := crypt.DecodeBase64(c.SignKey)
		encNumber, err := crypt.Encrypt([]byte(card.Number), sKey)
		if err != nil {
			return err
		}

		encDate, err := crypt.Encrypt([]byte(card.Date), sKey)
		if err != nil {
			return err
		}

		encCvv, err := crypt.Encrypt([]byte(card.Cvv), sKey)
		if err != nil {
			return err
		}

		encMeta, err := crypt.Encrypt([]byte(card.Meta), sKey)
		if err != nil {
			return err
		}

		card.Number = crypt.EncodeBase64(encNumber)
		card.Date = crypt.EncodeBase64(encDate)
		card.Cvv = crypt.EncodeBase64(encCvv)
		card.Meta = crypt.EncodeBase64(encMeta)
	}

	err = a.db.AddCard(card)
	if err != nil {
		return err
	}

	return err
}

func (a *App) GetCard(localID int) (card model.DataCard, err error) {

	c, err := a.GetUserConfig()
	if err != nil {
		return card, err
	}

	card, err = a.db.GetCard(localID)
	if err != nil {
		return card, err
	}

	sKey := crypt.DecodeBase64(c.SignKey)

	decNumber, err := crypt.Decrypt(crypt.DecodeBase64(card.Number), sKey)
	if err != nil {
		return card, err
	}

	decDate, err := crypt.Decrypt(crypt.DecodeBase64(card.Date), sKey)
	if err != nil {
		return card, err
	}

	decCvv, err := crypt.Decrypt(crypt.DecodeBase64(card.Cvv), sKey)
	if err != nil {
		return card, err
	}

	decMeta, err := crypt.Decrypt(crypt.DecodeBase64(card.Meta), sKey)
	if err != nil {
		return card, err
	}

	card.Number = string(decNumber)
	card.Date = string(decDate)
	card.Cvv = string(decCvv)
	card.Meta = string(decMeta)

	return card, err
}

func (a *App) GetAllCards() (cards []model.DataCard, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return cards, err
	}
	sKey := crypt.DecodeBase64(c.SignKey)

	cardEnc, err := a.db.GetAllCards()

	for _, v := range cardEnc {
		decNumber, err := crypt.Decrypt(crypt.DecodeBase64(v.Number), sKey)
		if err != nil {
			return cards, err
		}

		decDate, err := crypt.Decrypt(crypt.DecodeBase64(v.Date), sKey)
		if err != nil {
			return cards, err
		}

		decCvv, err := crypt.Decrypt(crypt.DecodeBase64(v.Cvv), sKey)
		if err != nil {
			return cards, err
		}

		decMeta, err := crypt.Decrypt(crypt.DecodeBase64(v.Meta), sKey)
		if err != nil {
			return cards, err
		}

		v.Number = string(decNumber)
		v.Date = string(decDate)
		v.Cvv = string(decCvv)
		v.Meta = string(decMeta)
		cards = append(cards, v)
	}

	return cards, err
}

func (a *App) DeleteCard(localID int) (err error) {

	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	item, err := a.GetCard(localID)
	if err != nil {
		return err
	}

	err = a.db.DeleteCard(localID)
	if err != nil {
		return err
	}

	go func() {
		err = a.HTTPService.DeleteCard(c.AccessToken, item.ExternalID)
		if err != nil {
			logger.Error("goroutine delete:", err)
		}
	}()

	return err
}

func (a *App) AddCred(cred *model.DataCred, encrypted bool) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	if !encrypted {
		sKey := crypt.DecodeBase64(c.SignKey)
		encUsername, err := crypt.Encrypt([]byte(cred.Username), sKey)
		if err != nil {
			return err
		}

		encPass, err := crypt.Encrypt([]byte(cred.Password), sKey)
		if err != nil {
			return err
		}

		encMeta, err := crypt.Encrypt([]byte(cred.Meta), sKey)
		if err != nil {
			return err
		}

		cred.Username = crypt.EncodeBase64(encUsername)
		cred.Password = crypt.EncodeBase64(encPass)
		cred.Meta = crypt.EncodeBase64(encMeta)
	}

	err = a.db.AddCred(cred)
	if err != nil {
		return err
	}

	return err
}

func (a *App) GetCred(localID int) (cred model.DataCred, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return cred, err
	}

	cred, err = a.db.GetCred(localID)
	if err != nil {
		return cred, err
	}

	sKey := crypt.DecodeBase64(c.SignKey)

	decUsername, err := crypt.Decrypt(crypt.DecodeBase64(cred.Username), sKey)
	if err != nil {
		return cred, err
	}

	decPass, err := crypt.Decrypt(crypt.DecodeBase64(cred.Password), sKey)
	if err != nil {
		return cred, err
	}

	decMeta, err := crypt.Decrypt(crypt.DecodeBase64(cred.Meta), sKey)
	if err != nil {
		return cred, err
	}

	cred.Username = string(decUsername)
	cred.Password = string(decPass)
	cred.Meta = string(decMeta)

	return cred, err
}

func (a *App) GetAllCreds() (creds []model.DataCred, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return creds, err
	}
	sKey := crypt.DecodeBase64(c.SignKey)

	credsEnc, err := a.db.GetAllCreds()

	for _, v := range credsEnc {

		decUsername, err := crypt.Decrypt(crypt.DecodeBase64(v.Username), sKey)
		if err != nil {
			return creds, err
		}

		decPass, err := crypt.Decrypt(crypt.DecodeBase64(v.Password), sKey)
		if err != nil {
			return creds, err
		}

		decMeta, err := crypt.Decrypt(crypt.DecodeBase64(v.Meta), sKey)
		if err != nil {
			return creds, err
		}

		v.Username = string(decUsername)
		v.Password = string(decPass)
		v.Meta = string(decMeta)

		creds = append(creds, v)
	}

	return creds, err
}

func (a *App) DeleteCred(localID int) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	item, err := a.GetCred(localID)
	if err != nil {
		return err
	}

	err = a.db.DeleteCred(localID)
	if err != nil {
		return err
	}

	go func() {
		err = a.HTTPService.DeleteCred(c.AccessToken, item.ExternalID)
		if err != nil {
			logger.Error("goroutine delete:", err)
		}
	}()

	return err
}

func (a *App) AddText(text *model.DataText, encrypted bool) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	if !encrypted {
		sKey := crypt.DecodeBase64(c.SignKey)
		encText, err := crypt.Encrypt([]byte(text.Text), sKey)
		if err != nil {
			return err
		}

		encMeta, err := crypt.Encrypt([]byte(text.Meta), sKey)
		if err != nil {
			return err
		}

		text.Text = crypt.EncodeBase64(encText)
		text.Meta = crypt.EncodeBase64(encMeta)
	}

	err = a.db.AddText(text)
	if err != nil {
		return err
	}

	return err
}

func (a *App) GetText(localID int) (text model.DataText, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return text, err
	}

	text, err = a.db.GetText(localID)
	if err != nil {
		return text, err
	}

	sKey := crypt.DecodeBase64(c.SignKey)

	decText, err := crypt.Decrypt(crypt.DecodeBase64(text.Text), sKey)
	if err != nil {
		return text, err
	}

	decMeta, err := crypt.Decrypt(crypt.DecodeBase64(text.Meta), sKey)
	if err != nil {
		return text, err
	}

	text.Text = string(decText)
	text.Meta = string(decMeta)

	return text, err
}

func (a *App) GetAllTexts() (texts []model.DataText, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return texts, err
	}
	sKey := crypt.DecodeBase64(c.SignKey)

	textsEnc, err := a.db.GetAllTexts()

	for _, v := range textsEnc {

		decText, err := crypt.Decrypt(crypt.DecodeBase64(v.Text), sKey)
		if err != nil {
			return texts, err
		}

		decMeta, err := crypt.Decrypt(crypt.DecodeBase64(v.Meta), sKey)
		if err != nil {
			return texts, err
		}

		v.Text = string(decText)
		v.Meta = string(decMeta)

		texts = append(texts, v)
	}

	return texts, err
}

func (a *App) DeleteText(localID int) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	item, err := a.GetText(localID)
	if err != nil {
		return err
	}

	err = a.db.DeleteText(localID)
	if err != nil {
		return err
	}

	go func() {
		err = a.HTTPService.DeleteText(c.AccessToken, item.ExternalID)
		if err != nil {
			logger.Error("goroutine delete:", err)
		}
	}()
	return err
}

func (a *App) AddFile(file *model.DataFile, encrypted bool) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	if !encrypted {
		sKey := crypt.DecodeBase64(c.SignKey)

		encMeta, err := crypt.Encrypt([]byte(file.Meta), sKey)
		if err != nil {
			return err
		}

		file.Meta = crypt.EncodeBase64(encMeta)
	}

	err = a.db.AddFile(file)
	if err != nil {
		return err
	}

	return err
}

func (a *App) GetFile(localID int) (file model.DataFile, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return file, err
	}

	file, err = a.db.GetFile(localID)
	if err != nil {
		return file, err
	}

	sKey := crypt.DecodeBase64(c.SignKey)

	decMeta, err := crypt.Decrypt(crypt.DecodeBase64(file.Meta), sKey)
	if err != nil {
		return file, err
	}

	file.Meta = string(decMeta)

	return file, err
}

func (a *App) GetAllFiles() (files []model.DataFile, err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return files, err
	}
	sKey := crypt.DecodeBase64(c.SignKey)

	filesEnc, err := a.db.GetAllFiles()

	for _, v := range filesEnc {

		decMeta, err := crypt.Decrypt(crypt.DecodeBase64(v.Meta), sKey)
		if err != nil {
			return files, err
		}

		v.Meta = string(decMeta)

		files = append(files, v)
	}

	return files, err
}

func (a *App) DeleteFile(localID int) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	item, err := a.GetFile(localID)
	if err != nil {
		return err
	}

	err = a.db.DeleteFile(localID)
	if err != nil {
		return err
	}

	err = a.FileService.DeleteFile(item.Path)
	if err != nil {
		return err
	}

	go func() {
		err = a.HTTPService.DeleteFile(c.AccessToken, item.ExternalID)
		if err != nil {
			logger.Error("goroutine delete:", err)
		}
	}()
	return err
}

func (a *App) SyncData() {
	c, err := a.GetUserConfig()
	if err != nil {
		dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)

		return
	}

	tokens, err := a.HTTPService.PostRefreshToken(c.RefreshToken)
	if err != nil {
		logger.Error(err)

		if err == service.ErrStatusUnauthorized {
			dialog.ShowError(errors.New("сессия устарела, авторизуйтесь повторно"), a.window)
			a.pageAuth()

			return
		}

		if err == service.ErrStatusUnauthorized {
			dialog.ShowError(errors.New("ошибка соединения с сервером"), a.window)

			return
		}
	}

	c.RefreshToken = tokens.RefreshToken
	c.AccessToken = tokens.AccessToken

	err = a.SetUserConfig(c)
	if err != nil {
		logger.Error(err)
		dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)
		return
	}

	err = a.SyncCards(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.Channels.SyncProgressBar <- 0.25

	err = a.SyncCreds(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.Channels.SyncProgressBar <- 0.50

	err = a.SyncTexts(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.Channels.SyncProgressBar <- 0.75

	err = a.SyncFiles(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.Channels.SyncProgressBar <- 1.0

	a.Channels.SyncProgressBarQuit <- true
}

func (a *App) SyncCards(accessToken string) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}
	cardsMap := make(map[int]*model.DataCard)
	getCardsMap := make(map[int]*model.DataCard)

	cards, err := a.db.GetAllCards()
	if err != nil {
		return err
	}

	for _, v := range cards {
		cardsMap[v.ExternalID] = &v
	}

	getCards, err := a.HTTPService.GetCardList(accessToken)
	if err != nil {
		return err
	}

	for _, v := range getCards {
		getCardsMap[v.ExternalID] = v
	}

	// создаем записи в бд клиента
	for _, v := range getCards {

		val, ok := cardsMap[v.ExternalID]
		if ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		updateCard := v
		if ok {
			updateCard.LocalID = val.LocalID
		}

		errUpdate := a.AddCard(updateCard, true)
		if errUpdate != nil {
			logger.Error("SyncCards - errUpdate local: ", errUpdate)
		}

	}

	// создаем записи в бд сервера
	for _, v := range cards {

		if val, ok := getCardsMap[v.ExternalID]; ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		// отправим на сервер
		reqBody := smodel.DataCard{
			ID:        v.ExternalID,
			Title:     v.Title,
			Number:    v.Number,
			Date:      v.Date,
			Cvv:       v.Cvv,
			Meta:      v.Meta,
			UpdatedAt: v.UpdatedAt,
		}

		id, err := a.HTTPService.AddCard(c.AccessToken, reqBody)
		if err != nil {
			logger.Error("goroutine add:", err)
		}
		v.ExternalID = id
		err = a.db.AddCard(&v)
		if err != nil {
			logger.Error("err --  ", err)
			//return
		}
	}

	return nil
}

func (a *App) SyncCreds(accessToken string) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}
	credsMap := make(map[int]*model.DataCred)
	getCredsMap := make(map[int]*model.DataCred)

	creds, err := a.db.GetAllCreds()
	if err != nil {
		return err
	}

	for _, v := range creds {
		credsMap[v.ExternalID] = &v
	}

	getCreds, err := a.HTTPService.GetCredList(accessToken)
	if err != nil {
		return err
	}

	for _, v := range getCreds {
		getCredsMap[v.ExternalID] = v
	}

	// создаем записи в бд клиента
	for _, v := range getCreds {
		if val, ok := credsMap[v.ExternalID]; ok {
			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				updateCred := v
				updateCred.LocalID = val.LocalID
				errUpdate := a.AddCred(updateCred, true)
				if errUpdate != nil {
					logger.Error(errUpdate)
				}
			}
		} else {
			newCred := model.DataCred{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Username:   v.Username,
				Password:   v.Password,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.AddCred(&newCred, true)
			if errAdd != nil {
				logger.Error(errAdd)
			}

		}
	}

	// создаем записи в бд сервера
	for _, v := range creds {

		if val, ok := getCredsMap[v.ExternalID]; ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		// отправим на сервер
		reqBody := smodel.DataCred{
			ID:        v.ExternalID,
			Title:     v.Title,
			Username:  v.Username,
			Password:  v.Password,
			Meta:      v.Meta,
			UpdatedAt: v.UpdatedAt,
		}

		id, err := a.HTTPService.AddCred(c.AccessToken, reqBody)
		if err != nil {
			logger.Error("goroutine add:", err)
		}
		v.ExternalID = id
		err = a.db.AddCred(&v)
		if err != nil {
			logger.Error("err --  ", err)
			//return
		}
	}

	return nil
}

func (a *App) SyncTexts(accessToken string) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}
	textsMap := make(map[int]*model.DataText)
	getTextsMap := make(map[int]*model.DataText)

	texts, err := a.db.GetAllTexts()
	if err != nil {
		return err
	}

	for _, v := range texts {
		textsMap[v.ExternalID] = &v
	}

	getTexts, err := a.HTTPService.GetTextList(accessToken)
	if err != nil {
		return err
	}

	for _, v := range getTexts {
		getTextsMap[v.ExternalID] = v
	}

	// создаем записи в бд клиента
	for _, v := range getTexts {
		if val, ok := textsMap[v.ExternalID]; ok {
			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				updateText := v
				updateText.LocalID = val.LocalID
				errUpdate := a.AddText(updateText, true)
				if errUpdate != nil {
					logger.Error(errUpdate)
				}
			}
		} else {
			newText := model.DataText{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Text:       v.Text,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.AddText(&newText, true)
			if errAdd != nil {
				logger.Error(errAdd)
			}

		}
	}

	// создаем записи в бд сервера
	for _, v := range texts {

		if val, ok := getTextsMap[v.ExternalID]; ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		// отправим на сервер
		reqBody := smodel.DataText{
			ID:        v.ExternalID,
			Title:     v.Title,
			Text:      v.Text,
			Meta:      v.Meta,
			UpdatedAt: v.UpdatedAt,
		}

		id, err := a.HTTPService.AddText(c.AccessToken, reqBody)
		if err != nil {
			logger.Error(err)
		}
		v.ExternalID = id
		err = a.db.AddText(&v)
		if err != nil {
			logger.Error("err --  ", err)
			//return
		}
	}

	return nil
}

func (a *App) SyncFiles(accessToken string) (err error) {
	c, err := a.GetUserConfig()
	if err != nil {
		return err
	}

	filesMap := make(map[int]*model.DataFile)
	getFilesMap := make(map[int]*model.DataFile)

	files, err := a.db.GetAllFiles()
	if err != nil {
		return err
	}

	for _, v := range files {
		log.Println("filesMap:", v.Title, v.LocalID, v.ExternalID)
		filesMap[v.ExternalID] = &v
	}

	getFiles, err := a.HTTPService.GetFileList(accessToken)
	if err != nil {
		return err
	}

	for _, v := range getFiles {
		log.Println("getFilesMap:", v.Title, v.ExternalID)
		getFilesMap[v.ExternalID] = v
	}

	// создаем записи в бд клиента
	for _, v := range getFiles {

		val, ok := filesMap[v.ExternalID]

		if ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		dFile, errDF := a.HTTPService.DownloadFile(v.Path)
		if errDF != nil {
			logger.Error("SyncFiles - downloadFile: ", v.Filename)
			continue
		}

		ext := filepath.Ext(v.Filename)

		filePath, errFP := a.FileService.SaveFile(dFile, ext)
		if errFP != nil {
			logger.Error("SyncFiles - fileService.SaveFile: ", errFP)
			continue
		}

		dFile.Close()

		updateFile := v
		if ok {
			updateFile.LocalID = val.LocalID
		}
		updateFile.Path = filePath
		errUpdate := a.AddFile(updateFile, true)
		if errUpdate != nil {
			logger.Error(errUpdate)
		}
	}

	// создаем записи в бд сервера
	for _, v := range files {

		log.Println("files ", v.LocalID, v.ExternalID)

		if val, ok := getFilesMap[v.ExternalID]; ok {
			if val.UpdatedAt.Unix() > v.UpdatedAt.Unix() {
				continue
			}
		}

		// отправим на сервер
		reqBody := smodel.DataFile{
			ID:        v.ExternalID,
			Title:     v.Title,
			Filename:  v.Filename,
			Path:      v.Path,
			Meta:      v.Meta,
			UpdatedAt: v.UpdatedAt,
		}

		id, err := a.HTTPService.AddFile(c.AccessToken, reqBody)
		if err != nil {
			logger.Error(err)
		}

		//log.Println("id", id)

		v.ExternalID = id
		err = a.db.AddFile(&v)
		if err != nil {
			logger.Error("err --  ", err)
			//return
		}
	}

	return nil
}
