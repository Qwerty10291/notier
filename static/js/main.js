// запрос всех записей. Возвращает список обьектов вида {id: 1, text: some text, created_at:1643198475 (это формат времени unix, количество секунд прошедших с 1 января 1970 года)}
async function LoadNotes(){
    try{
        let notes = await (await fetch("/notes")).json()
        console.log("notes loaded:", notes)
        return notes
    } catch (error) {
        console.log("failed to load notes:", error)
    }
}

// создание заметки. Возвращает такой же обьект как LoadNotes
async function CreateNote(text) {
    // создаем тело запроса
    let formData = new FormData()
    // поле text с текстом запроса
    formData.append("text", text)
    try{
        let response  = await fetch("/notes", {
            method: "POST",
            body: formData
        })
        let notes = await response.json()
        console.log("created note", notes)
        return notes
    } catch (error){
        console.log("error when creating note:", error)
    }
}

//удаление записи. Возвращает обьект вида: {status: false/true, error: ""}, status = true - запрос успешен, иначе в error будет текст ошибки
async function DeleteNote(noteId) {
    try{
        let response = await fetch("/notes/" + noteId, {
            method: "DELETE"
        })
        let deleteStatus = await response.json()
        console.log("delete " + noteId + "status:", deleteStatus)
        return deleteStatus
    } catch (error){
        console.log("error when deleting notes:", error);
    }
} 