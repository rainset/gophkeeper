package app

import (
	"context"
	"errors"
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
	"github.com/rainset/gophkeeper/internal/client/service/file"
	"github.com/rainset/gophkeeper/internal/client/storage"
	"github.com/rainset/gophkeeper/internal/client/ui"
	"github.com/rainset/gophkeeper/pkg/hash"
	"github.com/rainset/gophkeeper/pkg/logger"
	"image/color"
	"log"
	"net/url"
	"path/filepath"
	"time"
)

type App struct {
	window      fyne.Window
	cfg         *config.Config
	db          *storage.Base
	httpService *service.HttpService
	fileService *file.StorageFiles
	channels    channel.Channels
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

	fileService, err := file.New(filepath.Join(cfg.ClientFolder, "upload"))
	if err != nil {
		log.Fatal(err)
	}
	httpService := service.NewHttpService(cfg)

	return &App{
		window:      w,
		cfg:         cfg,
		db:          db,
		httpService: httpService,
		fileService: fileService,
		channels: channel.Channels{
			SyncProgressBar:     make(chan float64, 1),
			SyncProgressBarQuit: make(chan bool, 1),
		},
	}
}

func (a *App) Run() {

	a.PageAuth()
	//a.PageMain(ui.TypeCard)

	a.window.ShowAndRun()
	//a.PageAuth()
	//a.syncData()
	//a.PageMain(ui.TypeCard)

}

func (a *App) PageAuth() {

	authForm := a.authForm()
	regForm := a.regForm()

	tabs := container.NewAppTabs(
		container.NewTabItem("Авторизация", container.New(layout.NewPaddedLayout(), authForm)),
		container.NewTabItem("Регистрация", container.New(layout.NewPaddedLayout(), regForm)),
	)

	a.window.SetContent(tabs)
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

		a.db.SetUser(login.Text)

		c, err := a.db.GetConfig()
		if err != nil {
			dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)
			return
		}

		tokens, err := a.httpService.SignIn(model.User{Login: login.Text, Password: pass.Text})
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

		err = a.db.SetConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)
			return
		}

		a.PageMain(ui.TypeCard)

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
		tokens, err := a.httpService.SignUp(model.User{Login: login.Text, Password: pass.Text})
		if err != nil {

			logger.Error(err)

			if errors.Is(err, service.ErrStatusLoginExists) {
				dialog.ShowError(service.ErrStatusLoginExists, a.window)
				return
			}

			dialog.ShowError(service.ErrServer, a.window)
			return
		}

		a.db.SetUser(login.Text)

		c, err := a.db.GetConfig()
		if err != nil {
			dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)
			return
		}

		c.Login = login.Text
		c.Password = hash.Sha256(pass.Text)
		c.RefreshToken = tokens.RefreshToken
		c.AccessToken = tokens.AccessToken

		logger.Error(c)

		err = a.db.SetConfig(c)
		if err != nil {
			dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)
			return
		}

		a.PageMain(ui.TypeCard)
	}
	return regForm
}

func (a *App) PageMain(dataType ui.DataType) {

	var content *fyne.Container
	var tabs *container.AppTabs
	var syncBar *widget.ProgressBar

	syncBar = widget.NewProgressBar()
	syncBar.Hide()

	go func() {
		for {
			select {
			case v := <-a.channels.SyncProgressBar:
				syncBar.SetValue(v)
				logger.Info(v)
			case quit := <-a.channels.SyncProgressBarQuit:
				if quit {
					syncBar.Hide()
					tabs.Show()
					a.PageMain(ui.TypeCard)
				}
			}
		}
	}()

	addBtn := widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
		selectedTab := tabs.Selected()
		switch selectedTab.Text {
		case ui.TabCard.String():
			dataType = ui.TypeCard
			break
		case ui.TabCred.String():
			dataType = ui.TypeCred
			break
		case ui.TabText.String():
			dataType = ui.TypeText
			break
		case ui.TabFile.String():
			dataType = ui.TypeFile
			break
		}
		logger.Info("Добавить ", dataType)
		a.PageAdd(0, dataType)
		return
	})

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Синхронизация", theme.StorageIcon(), func() {
			syncBar.SetValue(0)
			syncBar.Show()
			tabs.Hide()
			a.syncData()
		}),
		addBtn,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Выйти", theme.ContentClearIcon(), func() {
			a.PageAuth()
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
		break
	case ui.TypeCred:
		tabs.Select(tabCred)
		break
	case ui.TypeText:
		tabs.Select(tabText)
		break
	case ui.TypeFile:
		tabs.Select(tabFile)
		break
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		syncBar,
		tabs,
	)

	a.window.SetContent(content)
}

func (a *App) PageAdd(ID int, dataType ui.DataType) {

	logger.Infof("pageadd %d %s", ID, dataType)

	var content *fyne.Container

	barTitle := canvas.NewText("Новая запись", color.Black)

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
			logger.Info("back button event ", dataType)
			a.PageMain(dataType)
		}),
		layout.NewSpacer(),
		barTitle,
	)

	tabCard := container.NewTabItem(ui.TabCard.String(), container.New(layout.NewPaddedLayout(), a.addCardForm(ID)))
	tabCred := container.NewTabItem(ui.TabCred.String(), container.New(layout.NewPaddedLayout(), a.addCredForm(ID)))
	tabText := container.NewTabItem(ui.TabText.String(), container.New(layout.NewPaddedLayout(), a.addTextForm(ID)))
	tabFile := container.NewTabItem(ui.TabFile.String(), container.New(layout.NewPaddedLayout(), a.addFileForm(ID)))

	tabs := container.NewAppTabs(
		tabCard,
		tabCred,
		tabText,
		tabFile,
	)

	switch dataType {
	case ui.TypeCard:
		tabs.Select(tabCard)
		break
	case ui.TypeCred:
		tabs.Select(tabCred)
		break
	case ui.TypeText:
		tabs.Select(tabText)
		break
	case ui.TypeFile:
		tabs.Select(tabFile)
		break
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		tabs,
	)

	a.window.SetContent(content)
}

func (a *App) PageEdit(ID int, dataType ui.DataType) {

	logger.Infof("PageEdit %d %s", ID, dataType)

	var content *fyne.Container

	//barTitle := canvas.NewText("Новая запись", color.Black)
	deleteBtn := widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {

		dialog.ShowConfirm("Удалить", "Вы действительно хотите удалить запись?", func(b bool) {
			if b {
				var err error
				c, err := a.db.GetConfig()
				if err != nil {
					dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)
					return
				}

				switch dataType {
				case ui.TypeCard:
					item, err := a.db.GetCard(ID)
					if err != nil {
						return
					}
					go a.httpService.DeleteCard(c.AccessToken, item.ExternalID)
					err = a.db.DeleteCard(ID)
					break
				case ui.TypeCred:
					item, err := a.db.GetCred(ID)
					if err != nil {
						return
					}
					go a.httpService.DeleteCred(c.AccessToken, item.ExternalID)
					err = a.db.DeleteCred(ID)
					break
				case ui.TypeText:
					item, err := a.db.GetText(ID)
					if err != nil {
						return
					}
					go a.httpService.DeleteText(c.AccessToken, item.ExternalID)
					err = a.db.DeleteText(ID)
					break
				case ui.TypeFile:
					item, err := a.db.GetFile(ID)
					if err != nil {
						return
					}
					go a.httpService.DeleteFile(c.AccessToken, item.ExternalID)
					err = a.db.DeleteFile(ID)
					break
				}
				if err != nil {
					dialog.ShowError(errors.New("ошибка при удалении записи"), a.window)
					return
				}
				a.PageMain(dataType)
			}

		}, a.window)

	})

	tasksBar := container.NewHBox(
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
			logger.Error("Назад ", dataType)
			a.PageMain(dataType)
		}),
		layout.NewSpacer(),
		deleteBtn,
	)

	var editCont *fyne.Container

	editCont = container.NewCenter()

	switch dataType {
	case ui.TypeCard:
		editCont = container.New(layout.NewPaddedLayout(), a.addCardForm(ID))
		break
	case ui.TypeCred:
		editCont = container.New(layout.NewPaddedLayout(), a.addCredForm(ID))
		break
	case ui.TypeText:
		editCont = container.New(layout.NewPaddedLayout(), a.addTextForm(ID))
		break
	case ui.TypeFile:
		editCont = container.New(layout.NewPaddedLayout(), a.addFileForm(ID))
		break
	}

	content = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		editCont,
	)

	a.window.SetContent(content)
}

func (a *App) addCardForm(ID int) *fyne.Container {

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

	if ID > 0 {
		item, err = a.db.GetCard(ID)
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
		a.PageMain(ui.TypeCard)
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

		if ID > 0 {
			cardData.LocalID = item.LocalID
		}

		err = a.db.AddCard(cardData)

		if err != nil {
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
			return
		}

		a.PageMain(ui.TypeCard)
	}
	return container.NewVBox(addForm)
}

func (a *App) addCredForm(ID int) *fyne.Container {

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

	if ID > 0 {
		item, err = a.db.GetCred(ID)
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
		a.PageMain(ui.TypeCred)
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

		if ID > 0 {
			credData.LocalID = item.LocalID
		}

		err = a.db.AddCred(credData)

		if err != nil {
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
			return
		}

		a.PageMain(ui.TypeCred)
	}
	return container.NewVBox(addForm)
}

func (a *App) addTextForm(ID int) *fyne.Container {

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

	if ID > 0 {
		item, err = a.db.GetText(ID)
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
		a.PageMain(ui.TypeText)
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

		if ID > 0 {
			textData.LocalID = item.LocalID
		}

		err = a.db.AddText(textData)

		if err != nil {
			dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
			return
		}

		a.PageMain(ui.TypeText)
	}
	return container.NewVBox(addForm)
}

func (a *App) addFileForm(ID int) *fyne.Container {

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

	if ID > 0 {
		item, err = a.db.GetFile(ID)
		if err != nil {
			logger.Error(err)
		}

		title.Text = item.Title
		uploadBtn.Text = item.Filename
		meta.Text = item.Meta
	}

	title.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")
	//text.Validator = validation.NewRegexp("^.{1,}", "обязательное поле")

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
		a.PageMain(ui.TypeFile)
	}
	addForm.SubmitText = "Сохранить"
	addForm.OnSubmit = func() {

		if tempFile.exists {
			filePath, err := a.fileService.SaveFile(tempFile.r, tempFile.ext)
			if err != nil {
				dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
				return
			}
			//filePathAbs, err := filepath.Abs(filePath)
			//if err != nil {
			//	dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
			//	return
			//}

			saveItem := model.DataFile{
				Title:     title.Text,
				Filename:  tempFile.filename,
				Path:      filePath,
				Ext:       tempFile.r.URI().Extension(),
				Meta:      meta.Text,
				UpdatedAt: time.Now(),
			}

			if ID > 0 {
				saveItem.LocalID = item.LocalID
			}

			err = a.db.AddFile(saveItem)
			if err != nil {
				logger.Error("a.db.AddFile: ", err)
				dialog.ShowError(errors.New("ошибка сохранения данных"), a.window)
				return
			}

			err = a.fileService.DeleteFile(item.Path)
			if err != nil {
				dialog.ShowError(errors.New("ошибка удаления старой версии файла"), a.window)
				return
			}
		}

		a.PageMain(ui.TypeFile)
	}

	if ID > 0 {

		fileDir := filepath.Dir(item.Path)
		u, err := url.Parse(fileDir)
		if err != nil {
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

	cards, err = a.db.GetAllCards()
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

		logger.Info("cardsList.OnSelected:", cards[id].LocalID, cards[id])
		logger.Info("cards[id].ID:", cards[id].LocalID)
		logger.Info("cards[id].Title:", cards[id].Title)
		logger.Info("cards[id].UpdatedAt:", cards[id].UpdatedAt)
		logger.Info("cards[id].ExternalID:", cards[id].ExternalID)
		logger.Info("cards[id].Number:", cards[id].Number)

		a.PageEdit(cards[id].LocalID, ui.TypeCard)
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

	creds, err = a.db.GetAllCreds()
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
		a.PageEdit(creds[id].LocalID, ui.TypeCred)
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

	texts, err = a.db.GetAllTexts()
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
		a.PageEdit(texts[id].LocalID, ui.TypeText)
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

	files, err = a.db.GetAllFiles()
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
		a.PageEdit(files[id].LocalID, ui.TypeFile)
	}

	scroll := container.NewScroll(filesList)
	scroll.SetMinSize(fyne.NewSize(100, 700))

	if len(files) != 0 {
		noItems.Hide()
	}
	return container.New(layout.NewPaddedLayout(), scroll, noItems)
}

func (a *App) syncData() {

	c, err := a.db.GetConfig()
	if err != nil {
		dialog.ShowError(errors.New("ошибка чтения настроек хранилища"), a.window)
		return
	}

	logger.Info(c)

	tokens, err := a.httpService.PostRefreshToken(c.RefreshToken)
	if err != nil {
		logger.Error(err)
		dialog.ShowError(errors.New("сессия устарела, авторизуйтесь повторно"), a.window)
		return
	}
	logger.Info("tokens ", tokens)

	c.RefreshToken = tokens.RefreshToken
	c.AccessToken = tokens.AccessToken

	err = a.db.SetConfig(c)
	if err != nil {
		logger.Error(err)
		dialog.ShowError(errors.New("ошибка записи настроек хранилища"), a.window)
		return
	}

	err = a.syncCards(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.channels.SyncProgressBar <- 0.25

	err = a.syncCreds(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.channels.SyncProgressBar <- 0.50

	err = a.syncTexts(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.channels.SyncProgressBar <- 0.75

	err = a.syncFiles(tokens.AccessToken)
	if err != nil {
		dialog.ShowError(errors.New("ошибка запроса списка с сервера"), a.window)
		return
	}

	a.channels.SyncProgressBar <- 1.0
	//time.Sleep(500 * time.Millisecond)
	a.channels.SyncProgressBarQuit <- true

}

func (a *App) syncCards(accessToken string) (err error) {

	cardsMap := make(map[int]model.DataCard)

	cards, err := a.db.GetAllCards()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range cards {
		cardsMap[v.ExternalID] = v
	}

	getCards, err := a.httpService.GetCardList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getCards {
		if val, ok := cardsMap[v.ExternalID]; ok {
			//update
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
				errUpdate := a.db.AddCard(updateCard)
				if errUpdate != nil {
					logger.Error("syncCards - errUpdate: ", errUpdate)
				}
				logger.Info("updateCard:", updateCard)
			}

		} else {
			//add

			newCard := model.DataCard{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Number:     v.Number,
				Date:       v.Date,
				Cvv:        v.Cvv,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.db.AddCard(newCard)
			if errAdd != nil {
				logger.Error("syncCards - errAdd: ", errAdd)
			}
			logger.Info("newCard:", newCard)
		}
	}

	logger.Info("getCards: ", getCards)
	return nil
}

func (a *App) syncCreds(accessToken string) (err error) {

	credsMap := make(map[int]model.DataCred)

	creds, err := a.db.GetAllCreds()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range creds {
		credsMap[v.ExternalID] = v
	}

	getCreds, err := a.httpService.GetCredList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getCreds {
		if val, ok := credsMap[v.ExternalID]; ok {
			//update
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
				errUpdate := a.db.AddCred(updateCred)
				if errUpdate != nil {
					logger.Error("syncCreds - errUpdate: ", errUpdate)
				}
				logger.Info("updateCred:", updateCred)
			}

		} else {
			//add

			newCred := model.DataCred{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Username:   v.Username,
				Password:   v.Password,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.db.AddCred(newCred)
			if errAdd != nil {
				logger.Error("syncCreds - errAdd: ", errAdd)
			}
			logger.Info("newCred:", newCred)
		}
	}

	logger.Info("getCreds: ", getCreds)
	return nil
}

func (a *App) syncTexts(accessToken string) (err error) {

	textsMap := make(map[int]model.DataText)

	texts, err := a.db.GetAllTexts()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range texts {
		textsMap[v.ExternalID] = v
	}

	getTexts, err := a.httpService.GetTextList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getTexts {
		if val, ok := textsMap[v.ExternalID]; ok {
			//update
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
				errUpdate := a.db.AddText(updateText)
				if errUpdate != nil {
					logger.Error("syncTexts - errUpdate: ", errUpdate)
				}
				logger.Info("updateCred:", updateText)
			}

		} else {
			//add

			newText := model.DataText{
				ExternalID: v.ExternalID,
				Title:      v.Title,
				Text:       v.Text,
				Meta:       v.Meta,
				UpdatedAt:  v.UpdatedAt,
			}

			errAdd := a.db.AddText(newText)
			if errAdd != nil {
				logger.Error("syncTexts - errAdd: ", errAdd)
			}
			logger.Info("newText:", newText)
		}
	}

	return nil
}

func (a *App) syncFiles(accessToken string) (err error) {

	filesMap := make(map[int]model.DataFile)

	files, err := a.db.GetAllFiles()
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, v := range files {
		filesMap[v.ExternalID] = v
	}

	getFiles, err := a.httpService.GetFileList(accessToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	// создаем записи в бд клиента
	for _, v := range getFiles {
		if val, ok := filesMap[v.ExternalID]; ok {
			//update
			// если дата на сервере новее обновим локальные данные
			if val.UpdatedAt.Unix() < v.UpdatedAt.Unix() {

				dFile, errDF := a.httpService.DownloadFile(v.Path)
				if errDF != nil {
					logger.Error("downloadFile: ", v.Filename)
					continue
				}

				ext := filepath.Ext(v.Filename)
				filePath, errFP := a.fileService.SaveFile(dFile, ext)
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
				errUpdate := a.db.AddFile(updateFile)
				if errUpdate != nil {
					logger.Error("syncFile - errUpdate: ", errUpdate)
				}
				logger.Info("updateFile:", updateFile)

				errD := a.fileService.DeleteFile(val.Path)
				if errD != nil {
					continue
				}

			}

		} else {
			//add

			dFile, errDF := a.httpService.DownloadFile(v.Path)
			if errDF != nil {
				logger.Error("downloadFile: ", v.Filename)
				continue
			}

			ext := filepath.Ext(v.Filename)
			filePath, errFP := a.fileService.SaveFile(dFile, ext)
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

			errAdd := a.db.AddFile(newFile)
			if errAdd != nil {
				logger.Error("syncFile - errAdd: ", errAdd)
			}
			logger.Info("newFile:", newFile)
		}
	}

	return nil
}
