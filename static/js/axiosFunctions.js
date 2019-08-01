
var repo = ""
var dir = ""
var fileSelected = false
function listRepos(clickedID) {
    fileSelected = false
    repo = ""
    console.log(clickedID)
    axios.post('http://127.0.0.1:8000', {
        operation: "listRepos",
        content: clickedID
    })
        .then(function (response) {
            console.log(response.data)
            var name = document.getElementById("files");
            document.getElementById("files").innerHTML = "";
            for (var I = 0; I < response.data.length; I++) {

                nameList = "<a id='" + response.data[I].Name + "' href='javascript:;' onclick='listFilesFromRepo(this.id);'> <span class='glyphicon glyphicon-folder-open'></span> " + response.data[I].Name + "</a>";
                //console.log(nameList)
                document.getElementById("files").innerHTML += nameList;
            }
            document.getElementById("markdown").value = ""
            document.getElementById("Create").disabled = true;
            document.getElementById("Update").disabled = true;
            document.getElementById("Delete").disabled = true;
        })
        .catch(function (error) {
            console.log(response)
        });
}


function listFilesFromDir(clickedID) {
    console.log(clickedID)

    if (clickedID == '...') {
        if (dir.length == 0) {
            listRepos('')
            return
        }
        dir = dir.split("/")
        dir.pop()
        if (fileSelected)
            dir.pop()
        if (dir.length == 0) {
            dir = ""
            listFilesFromRepo(repo)
            return
        }
        else
            dir = dir.join("/")
    }
    else
        dir = clickedID
    fileSelected = false

    axios.post('http://127.0.0.1:8000', {
        operation: 'listFilesFromDir',
        content: dir
    })
        .then(function (response) {
            console.log(response.data)
            var name = document.getElementById("files");
            document.getElementById("files").innerHTML = "";
            document.getElementById("files").innerHTML += "<a id='...' href='javascript:;' onclick='listFilesFromDir(this.id);'><span class='glyphicon glyphicon-arrow-left'></span></a>";
            for (var I = 0; I < response.data.length; I++) {
                if (response.data[I].Type == 'file')
                    nameList = "<a id='" + response.data[I].Path + "' href='javascript:;' onclick='getFileContent(this.id);'> <span class='glyphicon glyphicon-file'></span> " + response.data[I].Name + "</a>";
                else
                    nameList = "<a id='" + response.data[I].Path + "' href='javascript:;' onclick='listFilesFromDir(this.id);'> <span class='glyphicon glyphicon-folder-open'></span> " + response.data[I].Name + "</a>";
                document.getElementById("files").innerHTML += nameList;
            }
            document.getElementById("markdown").value = ""
            document.getElementById("Create").disabled = false;
            document.getElementById("Update").disabled = true;
            document.getElementById("Delete").disabled = true;
        })
        .catch(function (error) {
            console.log(response)
        });
}

function listFilesFromRepo(clickedID) {
    fileSelected = false
    console.log(clickedID)
    repo = clickedID
    axios.post('http://127.0.0.1:8000', {
        operation: 'listFilesFromRepo',
        content: clickedID
    })
        .then(function (response) {
            console.log(response.data)
            var name = document.getElementById("files");
            document.getElementById("files").innerHTML = "";
            document.getElementById("files").innerHTML += "<a id='...' href='javascript:;' onclick='listFilesFromDir(this.id);'><span class='glyphicon glyphicon-arrow-left'></span></a>";
            for (var I = 0; I < response.data.length; I++) {
                if (response.data[I].Type == 'file')
                    nameList = "<a id='" + response.data[I].Path + "' href='javascript:;' onclick='getFileContent(this.id);'> <span class='glyphicon glyphicon-file'></span> " + response.data[I].Name + "</a>";
                else
                    nameList = "<a id='" + response.data[I].Path + "' href='javascript:;' onclick='listFilesFromDir(this.id);'> <span class='glyphicon glyphicon-folder-open'></span> " + response.data[I].Name + "</a>";
                document.getElementById("files").innerHTML += nameList;
            }
            document.getElementById("markdown").value = ""
            document.getElementById("Create").disabled = false;
            document.getElementById("Update").disabled = true;
            document.getElementById("Delete").disabled = true;
        })
        .catch(function (error) {
            console.log(response)
        });
}

function getFileContent(clickedID) {
    fileSelected = true
    console.log("Repo: " + repo)
    console.log("getFileContent(" + clickedID + ")")
    dir = clickedID
    if (clickedID == '...') {
        listRepos(clickedID)
        document.getElementById("Update").disabled = true;
        document.getElementById("Delete").disabled = true;
        return
    }
    axios.post('http://127.0.0.1:8000', {
        operation: "getFileContent",
        content: clickedID,
        repo: repo
    })
        .then(function (response) {
            console.log(response.data)
            document.getElementById("markdown").value = response.data
            document.getElementById("Update").disabled = false;
            document.getElementById("Delete").disabled = false;
        })
        .catch(function (error) {
            console.log(response)
        });

}
