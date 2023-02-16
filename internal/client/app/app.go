package app

import (
	"context"
	"errors"
	"image/color"
	"log"
	"net/url"
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
	"github.com/rainset/gophkeeper/pkg/hash"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type App struct {
	window      fyne.Window
	cfg         *config.Config
	db          *storage.Base
	HTTPService *service.HTTPService
	FileService *service.FileService
	Channels    channel.Channels
}

func New(ctx context.Context, cfg *config.Config) *App {
	fyneApp := app.New()
	w := fyneApp.NewWindow("Gophkeeper")
	w.Resize(fyne.NewSize(500, 700))

	w.SetIcon(theme.FyneLogo())

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

	a.window.SetContent(tabs)
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
				logger.Info(v)
			case quit := <-a.Channels.SyncProgressBarQuit:
				if quit {
					syncBar.Hide()
					tabs.Show()
					a.pageMain(ui.TypeCard)
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
		logger.Info("Добавить ", dataType)
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
			logger.Info("back button event ", dataType)
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

		err = a.SetUserConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)

			return
		}

		a.pageMain(ui.TypeCard)
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

		err = a.SetUserConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)

			return
		}

		a.pageMain(ui.TypeCard)
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
		}

		err = a.AddCard(cardData)

		if err != nil {
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
		}

		err = a.AddCred(credData)

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
		}

		err = a.AddText(textData)

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
		if tempFile.exists {
			filePath, err := a.FileService.SaveFile(tempFile.r, tempFile.ext)
			if err != nil {
				dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)

				return
			}

			saveItem := model.DataFile{
				Title:     title.Text,
				Filename:  tempFile.filename,
				Path:      filePath,
				Ext:       tempFile.r.URI().Extension(),
				Meta:      meta.Text,
				UpdatedAt: time.Now(),
			}

			if localID > 0 {
				saveItem.LocalID = item.LocalID
			}

			err = a.AddFile(saveItem)
			if err != nil {
				logger.Error("AddFile: ", err)
				dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
				return
			}

			err = a.FileService.DeleteFile(item.Path)
			if err != nil {
				dialog.ShowError(errors.New("ошибка удаления старой версии файла"), a.window)
				return
			}
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

		return container.NewVBox(addForm, linkContainer, container.NewCenter(widget.NewLabel(item.Path)))
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

func (a *App) SyncData() {
	c, err := a.GetUserConfig()
	if err != nil {
		dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)

		return
	}

	logger.Info(c)

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

	logger.Info("tokens ", tokens)

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
	cardsMap := make(map[int]model.DataCard)

	cards, err := a.GetAllCards()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range cards {
		cardsMap[v.ExternalID] = v
	}

	getCards, err := a.HTTPService.GetCardList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getCards {
		if val, ok := cardsMap[v.ExternalID]; ok {

			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				updateCard := model.DataCard{
					LocalID:    val.LocalID,
					ExternalID: v.ExternalID,
					Title:      v.Title,
					Number:     v.Number,
					Date:       v.Date,
					Cvv:        v.Cvv,
					Meta:       v.Meta,
					UpdatedAt:  v.UpdatedAt,
				}
				errUpdate := a.AddCard(updateCard)
				if errUpdate != nil {
					logger.Error("SyncCards - errUpdate: ", errUpdate)
				}

				logger.Info("updateCard:", updateCard)
			}
		} else {
			newCard := model.DataCard{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Number:     v.Number,
				Date:       v.Date,
				Cvv:        v.Cvv,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.AddCard(newCard)
			if errAdd != nil {
				logger.Error("SyncCards - errAdd: ", errAdd)
			}

			logger.Info("newCard:", newCard)
		}
	}

	logger.Info("getCards: ", getCards)
	return nil
}

func (a *App) SyncCreds(accessToken string) (err error) {
	credsMap := make(map[int]model.DataCred)

	creds, err := a.GetAllCreds()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range creds {
		credsMap[v.ExternalID] = v
	}

	getCreds, err := a.HTTPService.GetCredList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getCreds {
		if val, ok := credsMap[v.ExternalID]; ok {

			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				updateCred := model.DataCred{
					LocalID:    val.LocalID,
					ExternalID: v.ExternalID,
					Title:      v.Title,
					Username:   v.Username,
					Password:   v.Password,
					Meta:       v.Meta,
					UpdatedAt:  v.UpdatedAt,
				}
				errUpdate := a.AddCred(updateCred)
				if errUpdate != nil {
					logger.Error("SyncCreds - errUpdate: ", errUpdate)
				}

				logger.Info("updateCred:", updateCred)
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

			errAdd := a.AddCred(newCred)
			if errAdd != nil {
				logger.Error("SyncCreds - errAdd: ", errAdd)
			}

			logger.Info("newCred:", newCred)
		}
	}

	logger.Info("getCreds: ", getCreds)
	return nil
}

func (a *App) SyncTexts(accessToken string) (err error) {
	textsMap := make(map[int]model.DataText)

	texts, err := a.GetAllTexts()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range texts {
		textsMap[v.ExternalID] = v
	}

	getTexts, err := a.HTTPService.GetTextList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getTexts {
		if val, ok := textsMap[v.ExternalID]; ok {

			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				updateText := model.DataText{
					LocalID:    val.LocalID,
					ExternalID: v.ExternalID,
					Title:      v.Title,
					Text:       v.Text,
					Meta:       v.Meta,
					UpdatedAt:  v.UpdatedAt,
				}
				errUpdate := a.AddText(updateText)
				if errUpdate != nil {
					logger.Error("SyncTexts - errUpdate: ", errUpdate)
				}

				logger.Info("updateCred:", updateText)
			}
		} else {
			newText := model.DataText{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Text:       v.Text,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.AddText(newText)
			if errAdd != nil {
				logger.Error("SyncTexts - errAdd: ", errAdd)
			}
			logger.Info("newText:", newText)
		}
	}

	return nil
}

func (a *App) SyncFiles(accessToken string) (err error) {
	filesMap := make(map[int]model.DataFile)

	files, err := a.GetAllFiles()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range files {
		filesMap[v.ExternalID] = v
	}

	getFiles, err := a.HTTPService.GetFileList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getFiles {
		if val, ok := filesMap[v.ExternalID]; ok {

			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {
				dFile, errDF := a.HTTPService.DownloadFile(v.Path)
				if errDF != nil {
					logger.Error("downloadFile: ", v.Filename)

					continue
				}

				ext := filepath.Ext(v.Filename)
				filePath, errFP := a.FileService.SaveFile(dFile, ext)
				if errFP != nil {
					logger.Error("fileService.SaveFile: ", errFP)

					continue
				}

				dFile.Close()

				updateFile := model.DataFile{
					LocalID:    val.LocalID,
					ExternalID: v.ExternalID,
					Title:      v.Title,
					Filename:   v.Filename,
					Path:       filePath,
					Meta:       v.Meta,
					UpdatedAt:  v.UpdatedAt,
				}
				errUpdate := a.AddFile(updateFile)
				if errUpdate != nil {
					logger.Error("syncFile - errUpdate: ", errUpdate)
				}

				logger.Info("updateFile:", updateFile)

				errD := a.FileService.DeleteFile(val.Path)
				if errD != nil {
					continue
				}
			}
		} else {
			dFile, errDF := a.HTTPService.DownloadFile(v.Path)
			if errDF != nil {
				logger.Error("downloadFile: ", v.Filename)

				continue
			}

			ext := filepath.Ext(v.Filename)
			filePath, errFP := a.FileService.SaveFile(dFile, ext)
			if errFP != nil {
				logger.Error("fileService.SaveFile: ", errFP)

				continue
			}

			dFile.Close()

			newFile := model.DataFile{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Filename:   v.Filename,
				Path:       filePath,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.AddFile(newFile)
			if errAdd != nil {
				logger.Error("syncFile - errAdd: ", errAdd)
			}
			logger.Info("newFile:", newFile)
		}
	}

	return nil
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

func (a *App) AddCard(card model.DataCard) (err error) {
	err = a.db.AddCard(card)
	return err
}

func (a *App) GetCard(localID int) (card model.DataCard, err error) {
	card, err = a.db.GetCard(localID)
	return card, err
}

func (a *App) GetAllCards() (cards []model.DataCard, err error) {
	cards, err = a.db.GetAllCards()
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

	err = a.db.DeleteCred(localID)
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

func (a *App) AddCred(cred model.DataCred) (err error) {
	err = a.db.AddCred(cred)
	return err
}

func (a *App) GetCred(localID int) (card model.DataCred, err error) {
	card, err = a.db.GetCred(localID)
	return card, err
}

func (a *App) GetAllCreds() (cards []model.DataCred, err error) {
	cards, err = a.db.GetAllCreds()
	return cards, err
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

func (a *App) AddText(text model.DataText) (err error) {
	err = a.db.AddText(text)
	return err
}

func (a *App) GetText(localID int) (text model.DataText, err error) {
	text, err = a.db.GetText(localID)
	return text, err
}

func (a *App) GetAllTexts() (texts []model.DataText, err error) {
	texts, err = a.db.GetAllTexts()
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

func (a *App) AddFile(file model.DataFile) (err error) {
	err = a.db.AddFile(file)
	return err
}

func (a *App) GetFile(localID int) (file model.DataFile, err error) {
	file, err = a.db.GetFile(localID)
	return file, err
}

func (a *App) GetAllFiles() (files []model.DataFile, err error) {
	files, err = a.db.GetAllFiles()
	return files, err
}

func (a *App) DeleteFile(localID int) (err error) {
	err = a.db.DeleteFile(localID)
	return err
}
