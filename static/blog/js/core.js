function Element(elemstr, isAbbr) {
    // console.log("elemstr", elemstr);
    if(isAbbr && elemstr) {
        let abbr = elemstr;
        let classNames = abbr.match(/(?<=\.)[0-9a-zA-Z-]+/g);
        let idvalue = abbr.match(/(?<=#)[a-zA-Z-]+/) && abbr.match(/(?<=#)[a-zA-Z-]+/)[0];
        let attrStart = abbr.indexOf('[');
        let attrEnd = abbr.lastIndexOf(']');
        let tagName = abbr.split(/[\.|#|\[]/)[0];
        this.element = document.createElement(tagName);

        if(classNames) {
            classNames.forEach(cls => {
                this.element.classList.add(cls)
            })
        }
        if(idvalue) {
            this.element.id = idvalue;
        }
        if(attrStart > 0 && attrEnd > 0) {
            var attrs = abbr.slice(attrStart+1, attrEnd);
            attrs = '{'+ attrs.replaceAll('=', ":") + '}';
            attrs = attrs.replaceAll(/(?<=:)[0-9a-zA-z-#_]+/g, '"$&",');
            attrs = attrs.replaceAll('" ', '", ');
            attrs = attrs.replaceAll(/[0-9a-zA-z-#_]+(?=:)/g, '"$&"');
            if(attrs.slice(-2, -1) === ',') {
                attrs = attrs.slice(0, -2)+'}'
            }
            let attrObj = JSON.parse(attrs)
            for(const key in attrObj) {
                this.element.setAttribute(key, attrObj[key])
            }
        }
    }
    else {
        this.element = elemstr;
    }
};

Element.prototype = {
    ready: function(fn) {
        document.addEventListener('DOMContentLoaded', fn);
    },
    addClass: function(...classNames) {
        this.element?.classList.add(...classNames);
        return this;
    },
    removeClass: function(...classNames) {
        console.log("...classsName:", ...classNames);
        this.element?.classList.remove(...classNames)
        return this;
    },
    hasClass: function(name) {
        return this.element?.classList.contains(name)
    },
    toggleClass: function(className) {
        this.element?.classList.toggle(className);
        return this;
    },
    addChild: function(childElement) {
        if(childElement.hasOwnProperty('element')) {
            this.element?.appendChild(childElement.element);
        }
        else {
            this.element?.appendChild(childElement);
        }
        return this;
    },
    addChildren: function(children) {
        console.log("children->", children);
        children.forEach(child=> {
            this.element?.appendChild(child);
        });
        return this;
    },
    addTo: function(parentElement) {
        if(parentElement.hasOwnProperty('element')) {
            parentElement.element.appendChild(this.element);
        }
        else {
            $(parentElement).appendChild(this.element);
        }
        return parentElement;
    },
    addNodes: function(domFragment) {
        this.element?.parentNode.insertBefore(domFragment, this.element.firstChild)
        return this;
    },
    addSibling: function(elem) {
        if(elem.hasOwnProperty('element')) {
            this.element?.insertAdjacentElement('afterend', elem.element);
        }
        else {
            this.element?.insertAdjacentElement('afterend', elem);
        }
        
        return this;
    },
    html: function(htmlstr) {
        this.element.innerHTML = htmlstr;
        return this;
    },
    text: function(textstring) {
        this.element.textContent = textstring;
        return this;
    },
    css: function(styles) {
        if(styles.includes(":")) {
            let splitStyles = styles.split(';');
            splitStyles.forEach( style => {
                if(!style) {
                    return;
                }
                let curStyle = style.split(':');
                let name = curStyle[0].trim();
                let val  = curStyle[1].trim();
                this.element.style[name] = val;
            })
        }else {
            return window.getComputedStyle(this.element, null).getPropertyValue(styles);
        }
        return this;
    },
    click: function(listener) {
        this.element?.addEventListener('click', listener);
        return this;
    },
    keyup: function(listener) {
        this.element?.addEventListener('keyup', listener);
        return this;
    },
    keydown: function(listener) {
        this.element?.addEventListener('keydown', listener);
        return this;
    },
};

const ProxyHandler = {
    get(target, property) {
        // console.log(`[get]\n\t- finding {${property}} in target:`, target);
        let element = target.element;
        let hasProperty = (property in target);
        // console.log(`\n\t- hasProperty: ${hasProperty}`);
        if(!hasProperty) {
            // console.log(`\n\t- NOT FOUND property: {${property}}`);
            if(typeof element[property] == 'function') {
                // console.log(`\n\t- Element[property]:${property}() is function`);
                let result = element[property].bind(element);
                // console.log("\n\t- result ->", result);
                return result;
            }
            // console.log("\n\t- get element[property]", element[property]);
            return element[property];
        }else {
            return target[property]
        }
    },
    set(target, property, value, receiver) {
        // console.log(`[set]\n\t- target:${target}\n\t- property:${property}`);
        // console.log("target:", target);
        if(property == 'addChild') {
            console.log("addChild called");
            console.log("value is ->", value);
            target.append(value);
            return true;
        }
        let hasProperty = (property in target);
        // console.log(`\n\t- hasProperty: ${hasProperty}`);
        if(!hasProperty && ('element' in target)) {
            target.element[property] = value;
            return true;
        }
    },
};

function ElementWrapper(selector) {

    const newElement = true;
    const wrapElement = false;
    try {
            if(selector.tagName) {
                return new Proxy(new Element(selector, wrapElement), ProxyHandler)
            }
            if(selector.charAt(0) === '@') {
                let abbrstr = selector.slice(1);
                return new Element(abbrstr, newElement)
            }
            let foundElement = document.querySelector(selector);
            if(foundElement) {
                return new Proxy(new Element(foundElement, wrapElement), ProxyHandler)
            }else {
                return new Element(selector, newElement)
            }
        }
    catch(e) {
        return new Element(selector, newElement) 
    }
};
let $ = ElementWrapper;

let Selection = {
    range: window.getSelection().getRangeAt(0).cloneRange(),
    addRange: (range) => {
        window.getSelection().addRange(range);
    },
    isRange: () => {
        if(window.getSelection().type === 'Range') {
            return true;
        }
        else {
            return false;
        }
    },
    removeAllRanges: () => {
        window.getSelection().removeAllRanges();
    },
    position: function() {
        if(window.getSelection().type === 'Range') {
            let rect = window.getSelection().getRangeAt(0).cloneRange().getBoundingClientRect();
            return rect;
        }
    },
    parentName: function() {
        if(this.sel.anchorNode.nodeType == 3) {
            return this.sel.anchorNode.parentElement.tagName;
        }
        else {
            return this.sel.anchorNode.tagName;
        }
    },
    wrapperNode: function() {
        var currentElement = window.getSelection().anchorNode.parentElement;
        while(currentElement && currentElement.tagName !== 'DIV') { 
            currentElement = currentElement.parentElement;
            if(currentElement && currentElement.tagName == 'DIV') {
                return currentElement
            }
        }
        return currentElement;
    },
    deletePairs: function() {
        console.log("this.sel:", this.sel);
        if(this.sel.anchorNode.nodeType == 3 && this.sel.anchorNode == '()') {
            this.sel.anchorNode.nodeValue = ''
        }
    },
};

let FAIcons = {
    'bold' :  {
        content: '<i class="fa fa-bold" aria-hidden="true"></i>',
        command: () => { document.execCommand('bold', false) }
    },
    'check':  {
        content: '<i class="fa fa-check"></i>',
    },
    'times':  {
        content: '<i class="fa fa-times"></i>',
    },
    'italic':  {
        content: '<i class="fa fa-italic" aria-hidden="true"></i>',
        command: () => { document.execCommand('italic', false) }
    },
    'underline':  {
        content: '<i class="fa fa-underline" aria-hidden="true"></i>',
        command: () => { document.execCommand('underline', false) }
    },
    'strikethrough':  {
        content: '<i class="fa fa-strikethrough" aria-hidden="true"></i>',
        command: () => { document.execCommand('strikeThrough', false) }
    },
    'superscript':  {
        content: '<i class="fa fa-superscript" aria-hidden="true"></i>',
        command: () => { document.execCommand('superscript', false) }
    },
    'subscript':  {
        content: '<i class="fa fa-subscript" aria-hidden="true"></i>',
        command: () => { document.execCommand('subscript', false) }
    },
    'link':  {
        content: '<i class="fa fa-link" aria-hidden="true"></i>',
    },
    'image':  {
        content: '<i class="fa fa-picture-o" aria-hidden="true"></i>',
    },
    'html':  {
        content: '<i class="fa fa-code" aria-hidden="true"></i>',
    },
    'table':  {
        content: '<i class="fa fa-table" aria-hidden="true"></i>',
    },
    'orderedlist':  {
        content: '<i class="fa fa-list-ol" aria-hidden="true"></i>',
        command: () => { document.execCommand('insertOrderedList', false) }
    },
    'unorderedlist':  {
        content: '<i class="fa fa-list-ul" aria-hidden="true"></i>',
        command: () => { document.execCommand('insertUnorderedList', false) }
    },
    'indent': {
        content: '<i class="fa fa-indent" aria-hidden="true"></i>',
        command: () => { document.execCommand('indent', false) }
    },
    'outdent': {
        content: '<i class="fa fa-outdent" aria-hidden="true"></i>',
        command: () => { document.execCommand('outdent', false) }
    },
    'justifyCenter': {
        content: '<i class="fa fa-align-center" aria-hidden="true"></i>',
        command: () => { document.execCommand('justifyCenter', false) }
    },
    'justifyFull': {
        content: '<i class="fa fa-align-justify" aria-hidden="true"></i>',
        command: () => { document.execCommand('justifyFull', false) }
    },
    'justifyLeft': {
        content: '<i class="fa fa-align-left" aria-hidden="true"></i>',
        command: () => { document.execCommand('justifyLeft', false) }
    },
    'justifyRight': {
        content: '<i class="fa fa-align-right" aria-hidden="true"></i>',
        command: () => { document.execCommand('justifyRight', false) }
    },
    'earser': {
        content: '<i class="fa fa-eraser" aria-hidden="true"></i>',
    },
    'quote': {
        content: '<i class="fa fa-quote-right" aria-hidden="true"></i>',
    },
    'code': {
        content: '<i class="fa fa-code fa-lg" aria-hidden="true"></i>',
    },
    'h1': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>1</sup>',
        command: () => { document.execCommand('heading', false, 'h1') }
    },
    'h2': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>2</sup>',
        command: () => { document.execCommand('heading', false, 'h2') }
    },
    'h3': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>3</sup>',
        command: () => { document.execCommand('heading', false, 'h3') }
    },
    'h4': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>4</sup>',
        command: () => { document.execCommand('heading', false, 'h4') }
    },
    'h5': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>5</sup>',
        command: () => { document.execCommand('heading', false, 'h5') }
    },
    'h6': {
        content: '<i class="fa fa-header" aria-hidden="true"><sup>6</sup>',
        command: () => { document.execCommand('heading', false, 'h6') }
    },
};

let WPIcons = {
    'code': { 
        content:'<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M20.8 10.7l-4.3-4.3-1.1 1.1 4.3 4.3c.1.1.1.3 0 .4l-4.3 4.3 1.1 1.1 4.3-4.3c.7-.8.7-1.9 0-2.6zM4.2 11.8l4.3-4.3-1-1-4.3 4.3c-.7.7-.7 1.8 0 2.5l4.3 4.3 1.1-1.1-4.3-4.3c-.2-.1-.2-.3-.1-.4z"></path></svg>',
        name: 'Code',
    },
    'image': {
        content:'<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 4.5h14c.3 0 .5.2.5.5v8.4l-3-2.9c-.3-.3-.8-.3-1 0L11.9 14 9 12c-.3-.2-.6-.2-.8 0l-3.6 2.6V5c-.1-.3.1-.5.4-.5zm14 15H5c-.3 0-.5-.2-.5-.5v-2.4l4.1-3 3 1.9c.3.2.7.2.9-.1L16 12l3.5 3.4V19c0 .3-.2.5-.5.5z"></path></svg>',
        name: 'Image',
    },
    'table': {
        content:'<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 4.5h14c.3 0 .5.2.5.5v3.5h-15V5c0-.3.2-.5.5-.5zm8 5.5h6.5v3.5H13V10zm-1.5 3.5h-7V10h7v3.5zm-7 5.5v-4h7v4.5H5c-.3 0-.5-.2-.5-.5zm14.5.5h-6V15h6.5v4c0 .3-.2.5-.5.5z"></path></svg>',
        name: 'Table',
    },
    'earser': {
        content:'<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-eraser" viewBox="0 0 24 24" stroke-width="1.5" stroke="#000000" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19 19h-11l-4 -4a1 1 0 0 1 0 -1.41l10 -10a1 1 0 0 1 1.41 0l5 5a1 1 0 0 1 0 1.41l-9 9" /><line x1="18" y1="12.3" x2="11.7" y2="6" /></svg>',
        name: 'unstyle',
    },
    'plus': {
        content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M18 11.2h-5.2V6h-1.6v5.2H6v1.6h5.2V18h1.6v-5.2H18z"></path></svg>',
    },
    'undo': {
        content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M18.3 11.7c-.6-.6-1.4-.9-2.3-.9H6.7l2.9-3.3-1.1-1-4.5 5L8.5 16l1-1-2.7-2.7H16c.5 0 .9.2 1.3.5 1 1 1 3.4 1 4.5v.3h1.5v-.2c0-1.5 0-4.3-1.5-5.7z"></path></svg>',
    },
    'redo': {
        content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M15.6 6.5l-1.1 1 2.9 3.3H8c-.9 0-1.7.3-2.3.9-1.4 1.5-1.4 4.2-1.4 5.6v.2h1.5v-.3c0-1.1 0-3.5 1-4.5.3-.3.7-.5 1.3-.5h9.2L14.5 15l1.1 1.1 4.6-4.6-4.6-5z"></path></svg>',
    },
    'setting': {
        content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path fill-rule="evenodd" d="M10.289 4.836A1 1 0 0111.275 4h1.306a1 1 0 01.987.836l.244 1.466c.787.26 1.503.679 2.108 1.218l1.393-.522a1 1 0 011.216.437l.653 1.13a1 1 0 01-.23 1.273l-1.148.944a6.025 6.025 0 010 2.435l1.149.946a1 1 0 01.23 1.272l-.653 1.13a1 1 0 01-1.216.437l-1.394-.522c-.605.54-1.32.958-2.108 1.218l-.244 1.466a1 1 0 01-.987.836h-1.306a1 1 0 01-.986-.836l-.244-1.466a5.995 5.995 0 01-2.108-1.218l-1.394.522a1 1 0 01-1.217-.436l-.653-1.131a1 1 0 01.23-1.272l1.149-.946a6.026 6.026 0 010-2.435l-1.148-.944a1 1 0 01-.23-1.272l.653-1.131a1 1 0 011.217-.437l1.393.522a5.994 5.994 0 012.108-1.218l.244-1.466zM14.929 12a3 3 0 11-6 0 3 3 0 016 0z" clip-rule="evenodd"></path></svg>',
    },
};

function createButtons(iconName) {
    
    this.buttonListener = {};
    if(iconName == 'wp') {
        this.buttonStock = WPIcons;
    }
    if(iconName == 'fa') {
        this.buttonStock = FAIcons;
    }
};
createButtons.prototype = {
    build: function(buttonNames) {
        // {buttNames}  : Array
        // {ButtonList} : DocumentFragment
        this.allButtons = new DocumentFragment();
        buttonNames.forEach( btn => {
            let button = $('@span')
                .html(this.buttonStock[btn].content)
                .addTo($('@button').addClass(btn));
            if(this.buttonStock[btn].name) {
                button.addChild($('@span').text(this.buttonStock[btn].name))
            }
            this.allButtons.appendChild(button.element)
            if(this.buttonStock[btn].command) {
                let execCommandFunc = this.buttonStock[btn].command;
                this.allButtons.querySelector('button.'+btn).addEventListener('click', function() {
                    execCommandFunc();
                });
            }
        });
        this.children = [...this.allButtons.children]
        return this;
    },
    buildwith: function(buttonNames, buttonWrapper) {
        // {buttNames}     : Array
        // {buttonWrapper} : String
        // {ButtonList}    : DocumentFragment
        this.allButtons = new DocumentFragment();
        buttonNames.forEach( btn => {
            let button = $('@span')
                .html(this.buttonStock[btn].content)
                .addTo($('@button').addClass(btn));
            let btw = document.createElement(buttonWrapper);
            btw.appendChild(button.element);
            if(this.buttonStock[btn].name) {
                button.addChild($('@span').text(this.buttonStock[btn].name))
            }
            this.allButtons.appendChild(btw);
            if(this.buttonStock[btn].command) {
                let execCommandFunc = this.buttonStock[btn].command;
                this.allButtons.querySelector('button.'+btn).addEventListener('click', function() {
                    execCommandFunc();
                });
            }
        });
        this.children = [...this.allButtons.children];
        return this;
    },
    wrapperElement: function(parentName) {
        this.wrapper = document.createElement(parentName);
        this.wrapper.appendChild(this.allButtons);
        return this;
    },
    wrapperClass: function(...classNames) {
        this.wrapper.classList.add(...classNames);
        return this;
    },
    button: function(buttonSelector) {
        if(this.wrapper?.childElementCount > 0) {
            return $(this.wrapper.querySelector(buttonSelector))
        }
        return $(this.allButtons.querySelector(buttonSelector));
    },
    getButton: function(btnName) {
        return this.buttonStock[btnName].content
    },
    forEach: function(callback) {
        this.children.forEach(child=> {
            callback($(child));
        });
    },
    getChildren: function(num) {
        if(num > 0) {
            return this.children.slice(0, num);
        }
        if(num < 0) {
            return this.children.slice(num);
        }
    },
    click: function(btn, listener) {
        if(this.buttonListener[btn] == null) {
            this.buttonListener[btn] = listener;
            this.button(btn).click(this.buttonListener[btn])
        }
        else {
            this.button(btn).click(this.buttonlistener[btn]);
        }
    },
    halt: function(btn) {
        this.button(btn).removeEventListener('click',this.buttonListener[btn]);
        this.buttonListener[btn]=null
        return this;
    },
};

function FontAwesomeButton() {}
FontAwesomeButton.prototype = new createButtons('fa');
FontAwesomeButton.prototype.New = function() {
    return new createButtons('fa');
};
let FAButton = new FontAwesomeButton();

function WordPressButton() {}
WordPressButton.prototype = new createButtons('wp');
WordPressButton.prototype.New = function() {
    return new createButtons('wp');
};
let WPButton = new WordPressButton();

export {
    ElementWrapper,
    createButtons,
    Selection,
    FAButton,
    WPButton,
};
