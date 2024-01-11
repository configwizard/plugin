import {RequestPlugins} from "../../wailsjs/go/main/Model";

const requestPlugins = async () => {
    try {
        const ids = await RequestPlugins();
        return ids; // Return the result
    } catch (e) {
        console.error("Error retrieving plugins: ", e);
        return []; // Return an empty array in case of an error
    }
};

export {
    requestPlugins
}
