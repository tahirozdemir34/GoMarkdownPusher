function listRepos(clickedID) {
    console.log(clickedID)
    axios.post('http://127.0.0.1:8000', {
        operation: "listRepos",
        content: clickedID
    })
    .then(function (response) {
        console.log(response.data)
        var name = document.getElementById("markdowns");
        document.getElementById("markdowns").innerHTML = "";
        for (var I = 0; I < response.data.length; I++) {
            nameList = "<a id='" + response.data[I].Name + "' href='javascript:;' onclick='listFilesFromRepo(this.id);'> <span class='glyphicon glyphicon-folder-open'></span> " + response.data[I].Name + "</a>";
            document.getElementById("markdowns").innerHTML += nameList;
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

var repo = ""
function listFilesFromRepo(clickedID) {
    console.log(clickedID)
    repo = clickedID
    axios.post('http://127.0.0.1:8000', {
        operation: "listFiles",
        content: clickedID
    })
    .then(function (response) {
        console.log(response.data)
        var name = document.getElementById("markdowns");
        document.getElementById("markdowns").innerHTML = "";
        for (var I = 0; I < response.data.length; I++) {

            nameList = "<a id='" + response.data[I].Name + "' href='javascript:;' onclick='getFileContent(this.id);'> <span class='glyphicon glyphicon-file'></span> " + response.data[I].Name + "</a>";
            document.getElementById("markdowns").innerHTML += nameList;
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
    console.log("Repo: " + repo)
    console.log("getFileContent(" + clickedID + ")")
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
