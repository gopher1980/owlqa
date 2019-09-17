function action() {
    console.log("Start open!!!")
    page.Navigate(param.Url)
    console.log("End Open!!!")
    return {message: "ok", action: action_name}
}