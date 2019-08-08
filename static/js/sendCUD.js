function sendCUD(clickedID) {

    var xhr = new XMLHttpRequest();
    xhr.open("POST", self.location.origin, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    if (clickedID == "Create") {
        if(document.getElementById("markdown").value.length == 0){
            alert("You must provide a content...")
            return
        }
        var fileName = prompt("Please enter file name", "");
        if (!isValidFileName(fileName)) {
            alert("You must enter a valid filename...")
        }
        else {
            if (fileSelected && dir.length>0) {
                dir = dir.split("/")
                dir.pop()
                dir = dir.join("/")
            }
            if (dir.length>0)
                dir += "/"
            xhr.send(JSON.stringify({
                operation: clickedID + "//" + dir+fileName,
                content: document.getElementById("markdown").value
            }));
            alert("File created...")
            console.log(fileSelected)
            if (dir.length == 0) {
                listFilesFromRepo(repo)
                dir += fileName
                getFileContent(dir)
                return
            }
            else {
                listFilesFromDir(dir)
                dir += fileName
                getFileContent(dir)
                return
            }

        }
    }
    else if (clickedID == "Update") {
        if(document.getElementById("markdown").value.length == 0){
            alert("You must provide a content...")
            return
        }
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
        if(dir.length == 0){
            listFilesFromRepo(repo)
            return
        }
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
        return rg1.test(fname) && !rg2.test(fname) && fname !=null;
    }
})();
