//@ts-check

function resetNewOpinionForm(){
    
    const form = /** @type {HTMLFormElement}*/ (document.getElementById("newOpinionForm"))

    form.reset()
}

function clearErrors(){
    const errorArea = document.getElementById("error-display")
    if(errorArea){
        errorArea.innerHTML = ""
    }
}
