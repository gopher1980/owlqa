function action() {
    console.log("Start click");
    var elem = Find(page, param.Strategy ,param.Selector)

    elem.Fill(param.Value);
    console.log("End click");
    return {message: "ok", action: action_name}
}