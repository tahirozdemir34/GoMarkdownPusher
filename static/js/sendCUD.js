function sendCUD(clickedID) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://127.0.0.1:8000", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    if (clickedID == "Create") {
        var fileName = prompt("Please enter file name", "");
        if (fileName != null) {
            fileName = fileName.replace(/\.[^/.]+$/, "")
            if(fileName.length == 0){
                alert("You must enter a valid filename...")
            }
            else{
                xhr.send(JSON.stringify({
                    operation: clickedID+"$"+fileName,
                    content: document.getElementById("markdown").value
                }));
            }
        }
        return
    }
    if (clickedID == "Update"){
        xhr.send(JSON.stringify({
            operation: clickedID,
            content: document.getElementById("markdown").value
        }));
        alert("File updated")
        return
    }
    if (confirm('Are you sure you want to delete this file?')) {
        xhr.send(JSON.stringify({
            operation: clickedID,
            content: "..."
        }));
        alert("File deleted")
        document.getElementById("markdown").value = ""
        location.reload()
    }
    

}