function sendCUD(clickedID) {

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://127.0.0.1:8000", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    if (clickedID == "Create") {
        var fileName = prompt("Please enter file name", "");
        if (!isValidFileName(fileName)) {
            alert("You must enter a valid filename...")
        }
        else {
            xhr.send(JSON.stringify({
                operation: clickedID + "/" + fileName,
                content: document.getElementById("markdown").value
            }));
            alert("File created...")

            if (dir.length == 0) {
                listFilesFromRepo(repo)
                return
            }
            else if (fileSelected) {
                dir = dir.split("/")
                dir.pop()
                dir = dir.join("/")
                listFilesFromDir(dir)
                dir += "/" + fileName
                getFileContent(dir)
                return
            }
            else {
                listFilesFromDir(dir)
                return
            }

        }
    }
    else if (clickedID == "Update") {

        xhr.send(JSON.stringify({
            operation: clickedID,
            content: document.getElementById("markdown").value
        }));
        alert("File updated")
    }
    else if (confirm('Are you sure you want to delete this file?')) {
        fileSelected = false
        xhr.send(JSON.stringify({
            operation: clickedID,
            content: "..."
        }));
        alert("File deleted")
        document.getElementById("markdown").value = ""
        dir = dir.split("/")
        dir.pop()
        dir = dir.join("/")
        listFilesFromDir(dir)
    }
}

var isValidFileName = (function () {
    var rg1 = /^[^\\/:\*\?"<>\|]+$/;
    var rg2 = /^(nul|prn|con|lpt[0-9]|com[0-9])(\.|$)/i;
    return function isValid(fname) {
        return rg1.test(fname) && !rg2.test(fname);
    }
})();