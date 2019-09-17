function action() {
    console.log("Start text");
    var elem = Find(page, param.Strategy ,param.Selector)
    str = elem.String()
    text = elem.Text()[0]
    console.log("End text");
    return {message: "ok", action: action_name,  string: str ,text: text}
}