function action() {
    console.log("Start sleep");
    sleep(param.Duration * Millisecond)
    console.log("End sleep");
    return {message: "ok", action: action_name}
}