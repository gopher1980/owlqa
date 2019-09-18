function click() {
    var elem = Find(page, param.Strategy, param.Selector, param.Frequency, param.Retry);
    elem.Click();
    return {message: "ok", action: "web.click"}
}

function fill() {
    var elem = Find(page, param.Strategy, param.Selector, param.Frequency, param.Retry);
    elem.Fill(param.Value);
    return {message: "ok", action: "web.fill"}
}

function navigate() {
    page.Navigate(param.Url);
    return {message: "ok", action: "web.navigate"}
}

function submit() {
    var elem = Find(page, param.Strategy, param.Selector, param.Frequency, param.Retry);
    elem.Submit();
    return {message: "ok", action: "web.submit"}
}

function text() {
    var elem = Find(page, param.Strategy, param.Selector, param.Frequency, param.Retry);
    value = elem.Text()[0];
    return {message: "ok", action: "web.text", text: value}
}

function waitShow() {
    var flag = WaitShowByStrategy(page, param.Strategy ,param.Selector, param.Frequency, param.Retry)
    return {message: "ok", action: "web.waitShow",  flag: flag}
}
function waitHide() {
    var flag = WaitHideByStrategy(page, param.Strategy ,param.Selector, param.Frequency, param.Retry)
    return {message: "ok", action: "web.waitHide",  flag: flag}
}