go build .
fyne package -os ios -appID com.sbarthelemy.networktools -icon Icon.png
fyne package -os android -appID com.sbarthelemy.networktools -icon Icon.png
#fyne package -os linux -icon myapp.png
#fyne package -os windows -icon myapp.png
fyne package -os darwin -icon Icon.png
