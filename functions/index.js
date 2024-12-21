const {onRequest} = require("firebase-functions/v2/https");
const logger = require("firebase-functions/logger");

// Example of a long line
const myLongFunctionName = someFunctionCallThatExceedsTheCharacterLimit().then(
    (response) => {
    // Code continues here...
    },
);

// Correct curly spacing
const obj = {name: "John"};
