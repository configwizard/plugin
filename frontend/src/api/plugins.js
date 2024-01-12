import {RequestPlugins} from "../../wailsjs/go/main/Model.js";

const requestPlugins = async () => {
    try {
        const plugins = await RequestPlugins();
        console.log("retrieved plugins ", plugins)
        return plugins; // Return the result
    } catch (e) {
        console.error("Error retrieving plugins: ", e);
        return []; // Return an empty array in case of an error
    }
};

export {
    requestPlugins
}
