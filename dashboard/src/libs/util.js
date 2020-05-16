let util = {

};
util.title = function (title) {
    title = title ? title + ' - Home' : 'Wizard Personal';
    window.document.title = title;
};

export default util;
