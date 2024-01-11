import {RequestContainers} from "../../wailsjs/go/main/Model";

const requestContainers = ()=> {
    try {
        RequestContainers().then(r => {})
    } catch(e) {
        console.log("error retrieving notification from database ", e)
    }
}

export {
    requestContainers
}
