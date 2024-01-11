
// Generates a random name for demonstration purposes
const generateRandomName = () => {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const randomIndex = Math.floor(Math.random() * characters.length);
    return characters[randomIndex] + Math.floor(Math.random() * 100);
};

// ... (previous code remains the same)

// Function to add a single entry
export const addSingleEntry = (array) => {
    return new Promise((resolve) => {
        const name = generateRandomName();
        const newArray = [...array, { Name: name }];
        newArray.sort((a, b) => a.Name.localeCompare(b.Name));
        setTimeout(() => resolve(newArray), 1000); // Resolves after 1 second
    });
};


// Function to add an entry to the array in alphabetical order
export const addEntryInAlphabeticalOrder = (array, maxEntries = 5) => {
    return new Promise((resolve) => {
        if (array.length < maxEntries) {
            const name = generateRandomName();
            array.push({ Name: name });
            array.sort((a, b) => a.Name.localeCompare(b.Name));
            setTimeout(() => resolve([...array]), 1000); // Resolves after 1 second
        } else {
            resolve([...array]);
        }
    });
};

