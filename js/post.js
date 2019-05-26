function postRequest(addr, params, callback) {
    req = new XMLHttpRequest();
    req.open("POST", addr, true);
    req.onreadystatechange = () => callback(req);
    req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    req.send(params);
}
