$(document).ready(function () {
    $('#tasktable').DataTable({
        responsive: true,
        "iDisplayLength": 5,
        "aLengthMenu": [[5, 10, 25], ["5 Per Page", "10 Per Page", "25 Per Page"]]
    });
});

const addGlobalEventListener = (
    type,
    selector,
    callback,
    options,
    parent = document
) => {
    parent.addEventListener(
        type,
        e => {
            if (e.target.matches(selector)) callback(e);
        },
        options
    );
}

addGlobalEventListener("click", ".resume-button", async function (e) {
    if (confirm("Do you wish to resume this task?") == true) {
        const el = e.target;
        const row = el.closest("tr");
        const resumeCell = row.querySelectorAll("td")[0];
        const resume = resumeCell.innerText;
        let payload = {
            taskid: resume,
            action: "resume"
        };
        console.log(payload);
        //send a post request with the data
        let res = await axios.post("/action", payload);
        let data = res.data;
        if (data == "done") {
            alert("TaskID " + resume + " Resumed.")
        }
    } else {
        alert("Action cancelled");
    }
}, {
    capture: true
});
addGlobalEventListener("click", ".pause-button", async function (e) {
    if (confirm("Do you wish to pause this task?") == true) {
        const el = e.target;
        const row = el.closest("tr");
        const pauseCell = row.querySelectorAll("td")[0];
        const pause = pauseCell.innerText;
        let payload = {
            taskid: pause,
            action: "pause"
        };
        console.log(payload);
        //send a post request with the data
        let res = await axios.post("/action", payload);
        let data = res.data;
        if (data == "done") {
            alert("TaskID " + pause + " Paused.")
        }
    } else {
        alert("Action cancelled");
    }
}, {
    capture: true
});

addGlobalEventListener("click", ".delete-button", async function (e) {
    if (confirm("Do you wish to delete this task?") == true) {
        const el = e.target;
        const row = el.closest("tr");
        const deleteCell = row.querySelectorAll("td")[0];
        const deletetask = deleteCell.innerText;
        let payload = {
            taskid: deletetask,
            action: "delete"
        };
        console.log(payload);
        //send a post request with the data
        let res = await axios.post("/action", payload);
        let data = res.data;
        if (data == "done") {
            alert("TaskID " + deletetask + " Deleted")
        }
    } else {
        alert("Action cancelled");
    }
}, {
    capture: true
});