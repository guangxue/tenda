
let BlogEditor = (function() {

    let Element = (function() {
        'use strict'
        function element(elemstr, isAbbr) {
            // console.log("elemstr", elemstr);
            if(isAbbr) {
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
        element.prototype = {
            ready: function(fn) {
                document.addEventListener('DOMContentLoaded', fn);
            },
            addClass: function(...classNames) {
                this.element.classList.add(...classNames);
                return this;
            },
            removeClass: function(...classNames) {
                this.element.classList.remove(...classNames)
                return this;
            },
            toggleClass: function(className) {
                this.element.classList.toggle(className);
                return this;
            },
            appendChild: function(childElement) {
                if(childElement.hasOwnProperty('element')) {
                    this.element.appendChild(childElement.element);
                }
                else {
                    this.element.appendChild(childElement);
                }
                return this;
            },

            appendTo: function(parentElement) {
                if(parentElement.hasOwnProperty('element')) {
                    parentElement.element.appendChild(this.element);
                }
                else {
                    $(parentElement).appendChild(this.element);
                }
                return parentElement;
            },
            firstChild: function(childElement) {
                this.element.insertAdjacentElement('afterbegin', childElement);
                return this;
            },
            lastChild: function(childElement) {
                this.element.insertAdjacentElement('beforeend', childElement);
                return this;
            },
            before: function(parentElement) {
                parentElement.element.insertAdjacentElement('beforebegin', this.element);
                return this;
            },
            after: function(parentElement) {
                parentElement.element.insertAdjacentElement('afterend', this.element);
                return this;
            },
            attachInside: function(liveElem) {
                if (typeof liveElem === 'string') {
                    $(liveElem).insertAdjacentElement('beforeend', this.element);
                }
                else {
                    liveElem.insertAdjacentElement('beforeend', this.element);
                }
                return this;
            },
            attachFragmentAfter: function(domFragment) {
                this.element.parentNode.insertBefore(domFragment, this.element.nextSibling)
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
                    this.element.style.cssText = styles;
                }else {
                    return window.getComputedStyle(this.element, null).getPropertyValue(styles);
                }
                return this;
            },
            click: function(listener) {
                this.element.addEventListener('click', listener);
                return this;
            },
            keyup: function(listener) {
                this.element.addEventListener('keyup', listener);
                return this;
            },
            keydown: function(listener) {
                this.element.addEventListener('keydown', listener);
                return this;
            },
        };
        return element;
    })();

    const ProxyHandler = {
        get(target, property) {
            // console.log(`[get]\n\t- finding {${property}} in target:`, target);
            let element = target.element;
            let hasProperty = (property in target);
            // console.log(`\n\t- hasProperty: ${hasProperty}`);
            if(property == 'addChild') {
                console.log("addChild called");
            }
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
    }

    let $ = function(selector) {
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
                if(selector.charAt(0) === '&') {
                    return new DocumentFragment();
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

    let Selection = (function() {
        'use strict'

        let selection = function() {
            this.selection = window.getSelection();
        };

        selection.prototype = {
            range: function() {
                return this.selection.getRangeAt(0).cloneRange();
            },
            parentName: function() {
                if(this.selection.anchorNode.nodeType == 3) {
                    return this.selection.anchorNode.parentElement.tagName;
                }
                else {
                    return this.selection.anchorNode.tagName;
                }
            },
            deletePairs: function() {
                console.log("this.selection:", this.selection);
                if(this.selection.anchorNode.nodeType == 3 && this.selection.anchorNode== '()') {
                    this.selection.anchorNode.nodeValue = ''
                }
            },
        };

        return new selection();
    })();

    // let TablerIcons = 
    let ToolbarButtons = (function() {
        'use strict'
        let toolbar = function(buttons) {
            this.buttonBank = {
                'plus': '<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-plus" width="20" height="20" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg>',
                'undo': '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M18.3 11.7c-.6-.6-1.4-.9-2.3-.9H6.7l2.9-3.3-1.1-1-4.5 5L8.5 16l1-1-2.7-2.7H16c.5 0 .9.2 1.3.5 1 1 1 3.4 1 4.5v.3h1.5v-.2c0-1.5 0-4.3-1.5-5.7z"></path></svg>',
                'redo': '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M15.6 6.5l-1.1 1 2.9 3.3H8c-.9 0-1.7.3-2.3.9-1.4 1.5-1.4 4.2-1.4 5.6v.2h1.5v-.3c0-1.1 0-3.5 1-4.5.3-.3.7-.5 1.3-.5h9.2L14.5 15l1.1 1.1 4.6-4.6-4.6-5z"></path></svg>',
                'settings': '<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-settings" width="24" height="24" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z" /><circle cx="12" cy="12" r="3" /></svg>',
            }
        };

        toolbar.prototype = {
            init: function(bulkButtons) {},
            createButtonList: function(button) {
                let buttonList = new DocumentFragment();
                let foundButtons = Object.keys(this.buttonBank).filter(bbtn => {
                    if(!buttons) {
                        return true;
                    }
                    else {
                        for(const btn of btnNames.values()) {
                            if(btn == bbtn) {
                                return true;
                            }
                        }
                    }
                });
                foundButtons.forEach( fbtn=> {
                    $('@button')
                        .appendChild($('@span').html())
                })
            },
            addButton: function() {},
            click: function(listener) {},
            button: function() {},
            replaceButton: function() {},
        };
        return toolbar;
    })();

    let navtoolbarBtns = {
        'plus-navbutton': {
            content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M18 11.2h-5.2V6h-1.6v5.2H6v1.6h5.2V18h1.6v-5.2H18z"></path></svg>',
            action: () => {},
        },
        'undo-navbutton': {
            content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M18.3 11.7c-.6-.6-1.4-.9-2.3-.9H6.7l2.9-3.3-1.1-1-4.5 5L8.5 16l1-1-2.7-2.7H16c.5 0 .9.2 1.3.5 1 1 1 3.4 1 4.5v.3h1.5v-.2c0-1.5 0-4.3-1.5-5.7z"></path></svg>',
            action: () => {},
        },
        'redo-navbutton': {
            content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path d="M15.6 6.5l-1.1 1 2.9 3.3H8c-.9 0-1.7.3-2.3.9-1.4 1.5-1.4 4.2-1.4 5.6v.2h1.5v-.3c0-1.1 0-3.5 1-4.5.3-.3.7-.5 1.3-.5h9.2L14.5 15l1.1 1.1 4.6-4.6-4.6-5z"></path></svg>',
            action: () => {},
        },
    };

    let navsettingBtns = {
        'preview': {
            content: 'Preview',
            action: () => {},
        },
        'publish': {
            content: 'Publish',
            action: () => {},
        },
        'settings': {
            content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true" focusable="false"><path fill-rule="evenodd" d="M10.289 4.836A1 1 0 0111.275 4h1.306a1 1 0 01.987.836l.244 1.466c.787.26 1.503.679 2.108 1.218l1.393-.522a1 1 0 011.216.437l.653 1.13a1 1 0 01-.23 1.273l-1.148.944a6.025 6.025 0 010 2.435l1.149.946a1 1 0 01.23 1.272l-.653 1.13a1 1 0 01-1.216.437l-1.394-.522c-.605.54-1.32.958-2.108 1.218l-.244 1.466a1 1 0 01-.987.836h-1.306a1 1 0 01-.986-.836l-.244-1.466a5.995 5.995 0 01-2.108-1.218l-1.394.522a1 1 0 01-1.217-.436l-.653-1.131a1 1 0 01.23-1.272l1.149-.946a6.026 6.026 0 010-2.435l-1.148-.944a1 1 0 01-.23-1.272l.653-1.131a1 1 0 011.217-.437l1.393.522a5.994 5.994 0 012.108-1.218l.244-1.466zM14.929 12a3 3 0 11-6 0 3 3 0 016 0z" clip-rule="evenodd"></path></svg>',
            action: () => {},
        }
    };

    let Navbar = (function() {
        'use strict'
        let navbar = function(navbuttons, settingbuttons) {
                $('header nav')
                    .appendChild(this.navtoolbar(navbuttons))
                    .appendChild($('@div.middleflex'))
                    .appendChild(this.navsetting(settingbuttons))
        };
        navbar.prototype = {
            navtoolbar: function(buttons) {
                for(const btn in buttons) {
                    let buttonList = new DocumentFragment();
                    Object.keys(buttons).forEach(btn=>{
                        let buttonDIV = $('@button').addClass(btn)
                            .html(buttons[btn].content)
                            .appendTo($('@div')).element;
                        buttonList.appendChild(buttonDIV);
                    })
                    return $('@div.navtoolbar')
                            .appendChild($('@div')
                            .css("height:60px;width:60px;background: #000;color:#fff;margin-right:1.5em;")
                            .appendChild($('@a[href=#].dashboard').text('').css('color:#fff')))
                            .appendChild(buttonList);
                }
            },
            navsetting: function(buttons) {
                for(const btn in buttons) {
                    let buttonList = new DocumentFragment();
                    Object.keys(buttons).forEach(btn=>{
                        let buttonDIV = $('@button').addClass(btn)
                            .html(buttons[btn].content)
                            .appendTo($('@div')).element;
                        buttonList.appendChild(buttonDIV);
                    })
                    return $('@div.navsettings').appendChild(buttonList);
                }
            },
        };
        return navbar;
    })();

    let SidebarButtons = {
        'code': {
            text: 'Code',
            content: '<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M20.8 10.7l-4.3-4.3-1.1 1.1 4.3 4.3c.1.1.1.3 0 .4l-4.3 4.3 1.1 1.1 4.3-4.3c.7-.8.7-1.9 0-2.6zM4.2 11.8l4.3-4.3-1-1-4.3 4.3c-.7.7-.7 1.8 0 2.5l4.3 4.3 1.1-1.1-4.3-4.3c-.2-.1-.2-.3-.1-.4z"></path></svg>',
            action: () => {
                let codeWrapper = $('pre')
                    .css('border: 1px solid black;padding: 1em 1em;')
                    .appendChild($('code').css('padding:5px;'))
                    .element;
                Selection.range().insertNode(codeWrapper);
            },
        },
        'image': {
            text: 'Image',
            contentee: '<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-photo" width="26" height="26" viewBox="0 0 24 24" stroke-width="1.5" stroke="#2c3e50" fill="none" stroke-linecap="round" stroke-linejoin="round"> <path stroke="none" d="M0 0h24v24H0z" fill="none"/><line x1="15" y1="8" x2="15.01" y2="8" /><rect x="4" y="4" width="16" height="16" rx="3" /><path d="M4 15l4 -4a3 5 0 0 1 3 0l5 5" /><path d="M14 14l1 -1a3 5 0 0 1 3 0l2 2" /></svg>',
            content: '<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 4.5h14c.3 0 .5.2.5.5v8.4l-3-2.9c-.3-.3-.8-.3-1 0L11.9 14 9 12c-.3-.2-.6-.2-.8 0l-3.6 2.6V5c-.1-.3.1-.5.4-.5zm14 15H5c-.3 0-.5-.2-.5-.5v-2.4l4.1-3 3 1.9c.3.2.7.2.9-.1L16 12l3.5 3.4V19c0 .3-.2.5-.5.5z"></path></svg>',
            action: () => {},
        },
        'table': {
            text: 'Table',
            contentee: '<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-table" width="26" height="26" viewBox="0 0 24 24" stroke-width="1.5" stroke="#2c3e50" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><rect x="4" y="4" width="16" height="16" rx="2" /><line x1="4" y1="10" x2="20" y2="10" /><line x1="10" y1="4" x2="10" y2="20" /></svg>',
            content: '<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" role="img" aria-hidden="true" focusable="false"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 4.5h14c.3 0 .5.2.5.5v3.5h-15V5c0-.3.2-.5.5-.5zm8 5.5h6.5v3.5H13V10zm-1.5 3.5h-7V10h7v3.5zm-7 5.5v-4h7v4.5H5c-.3 0-.5-.2-.5-.5zm14.5.5h-6V15h6.5v4c0 .3-.2.5-.5.5z"></path></svg>',
            action: () => {},
        },
        'earser': {
            text: 'Remove format',
            content: '<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-eraser" viewBox="0 0 24 24" stroke-width="1.5" stroke="#000000" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19 19h-11l-4 -4a1 1 0 0 1 0 -1.41l10 -10a1 1 0 0 1 1.41 0l5 5a1 1 0 0 1 0 1.41l-9 9" /><line x1="18" y1="12.3" x2="11.7" y2="6" /></svg>',
            action: () => {
                console.log("[earser] removing format");

                let tagName = window.getSelection().focusNode.parentElement.tagName;
                console.log("tagName", tagName);
                if(tagName == 'DIV' || tagName == 'ARTICLE') {
                    return;
                }
                else {
                    window.getSelection().focusNode.parentElement.outerHTML = window.getSelection().focusNode.parentElement.textContent;
                }
            },
        },
    };

    let SelectionButtons = {
        'h1': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>1</sup>',
            action: function() { document.execCommand('heading', false, 'H1') },
        },
        'h2': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>2</sup>',
            action: function() { document.execCommand('heading', false, 'H2') },
        },
        'h3': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>3</sup>',
            action: function() { document.execCommand('heading', false, 'H3') },
        },
        'h4': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>4</sup>',
            action: function() { document.execCommand('heading', false, 'H4') },
        },
        'h5': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>5</sup>',
            action: function() { document.execCommand('heading', false, 'H5') },
        },
        'h6': {
            className: 'button.heading',
            content: '<i class="fa fa-header" aria-hidden="true"><sup>6</sup>',
            action: function() { document.execCommand('heading', false, 'H6') },
        },
        'bold': {
            className: 'button.bold',
            content: '<i class="fa fa-bold" aria-hidden="true"></i>',
            action: function() { document.execCommand('bold', false) },
        },
        'italic': {
            className: 'button.italic',
            content: '<i class="fa fa-italic" aria-hidden="true"></i>',
            action: function() { document.execCommand('italic', false) },
        },
        'underline': {
            className: 'button.underline',
            content: '<i class="fa fa-underline" aria-hidden="true"></i>',
            action: function() { document.execCommand('underline', false) },
        },
        'strikethrough': {
            className: 'button.strikeThrough',
            content: '<i class="fa fa-strikethrough" aria-hidden="true"></i>',
            action: function() { document.execCommand('strikeThrough', false)},
        },
        'superscript': {
            className: 'button.superscript',
            content: '<i class="fa fa-superscript" aria-hidden="true"></i>',
            action: function() { document.execCommand('superscript', false)},
        },
        'subscript': {
            className: 'button.subscript',
            content: '<i class="fa fa-subscript" aria-hidden="true"></i>',
            action: function() { document.execCommand('subscript', false)},
        },
        
        'orderedlist': {
            className: 'button.insertOrderedList',
            content: '<i class="fa fa-list-ol" aria-hidden="true"></i>',
            action: function() { document.execCommand('insertOrderedList', false)},
        },
        'unorderedlist': {
            className: 'button.insertUnorderedList',
            content: '<i class="fa fa-list-ul" aria-hidden="true"></i>',
            action: function() { document.execCommand('insertUnorderedList', false)},
        },
        'indent': {
            className: 'button.indent',
            content: '<i class="fa fa-indent" aria-hidden="true"></i>',
            action: function() { document.execCommand("indent", false) },
        },
        'outdent': {
            className: 'button.outdent',
            content: '<i class="fa fa-outdent" aria-hidden="true"></i>',
            action: () => { document.execCommand("outdent", false) },
        },
        'justifyCenter': {
            className: 'button.justifyCenter',
            content: '<i class="fa fa-align-center" aria-hidden="true"></i>',
            action: () => { document.execCommand('justifyCenter', false) },
        },
        'justifyFull': {
            className: 'button.justifyFull',
            content: '<i class="fa fa-align-justify" aria-hidden="true"></i>',
            action: () => { document.execCommand('justifyFull', false) },
        },
        'justifyLeft': {
            className: 'button.justifyLeft',
            content: '<i class="fa fa-align-left" aria-hidden="true"></i>',
            action: () => { document.execCommand('justifyLeft', false) },
        },
        'justifyRight': {
            className: 'button.justifyRight',
            content: '<i class="fa fa-align-right" aria-hidden="true"></i>',
            action: () => { document.execCommand('justifyRight', false) },
        },
        'link': {
            className: 'button.createLink',
            ValueArgs: 'http',
            placeholder: 'Enter a link',
            content: '<i class="fa fa-link" aria-hidden="true"></i>',
            action: (range) => {
                // $('.editor-toolbar').addClass('displaynone');
                $(".inputToolbar").removeClass('displaynone').addClass('displayflex');
                document.getSelection().removeAllRanges();
                $('.inputToolbar input').focus();
                document.getSelection().addRange(range)
                
                $('a.save').click(function() {
                    let url = $('input').value;
                    console.log("url ->", url);
                    document.execCommand('createLink', false, url);
                    $(".inputToolbar").removeClass('displayflex').addClass("displaynone");
                    $(".selectionToolbar").removeClass('displaynone');
                    $('.inputToolbar input').value = '';
                });

                $('a.close').click(function() {
                    $(".inputToolbar").removeClass('displayflex').addClass("displaynone");
                    $(".selectionToolbar").removeClass('displaynone');
                    $('.inputToolbar input').value = "";
                });
            },
        },
        'quote': {
            className: 'button.quote',
            content: '<i class="fa fa-quote-right" aria-hidden="true"></i>',
        },
        'removeFormat': {
            className: 'button.removeFormat',
            content: '<i class="fa fa-eraser" aria-hidden="true"></i>',
            action: (range) => {
                
            },
        },
        'image': {
            className: 'button.insertImage',
            ValueArgs: 'image',
            placeholder: 'Enter a image link',
            content: '<i class="fa fa-picture-o" aria-hidden="true"></i>'
        },
        'html': {
            className: 'button.code',
            content: '<i class="fa fa-code" aria-hidden="true"></i>'
        },
        'table': {
            className: 'button.table',
            content: '<i class="fa fa-table" aria-hidden="true"></i>',
        },    
    };

    let LeftSideToolbar = (function() {
        'use strict'

        let leftsidebar = function(buttons) {

                $('.left-menu-wrapper .textheader').attachFragmentAfter(this.createLeftPanelButtons(buttons))
        };
        leftsidebar.prototype = {
            createLeftPanelButtons: function(buttons) {
                if(!buttons) {
                    buttons = ['code', 'image', 'table', 'earser'];
                }
                let buttonList = new DocumentFragment();
                let foundButtons = Object.keys(SidebarButtons).filter(btn=> {
                    for(const value of buttons.values()) {
                        if(value == btn) {
                            return true;
                        }
                    }
                });
                foundButtons.forEach( btn=> {
                    let buttonDIV = $('@button')
                        .click(function() {
                            if(SidebarButtons[btn].action) {
                                SidebarButtons[btn].action();
                            }
                        })
                        .appendChild($('@span').html(SidebarButtons[btn].content))
                        .appendChild($('@span').text(SidebarButtons[btn].text))
                        .appendTo($('@div'))
                        .element;
                    buttonList.appendChild(buttonDIV)
                });
                return buttonList;
            },
        };
        return leftsidebar
    })();
    
    let SelectionToolbar = (function() {
        'use strict'

        let toolbar = function(defaultSettings) {
            this.location = defaultSettings.location;
            this.locElement = this.locElement
        };
        /*

            let leftToolbar = Toolbar('.left-menu-wrapper', SidebarButtons.default)
        */

        toolbar.prototype = {
            init: function(btns) {
                return this.createToolbar(btns);
            },
            createToolbar: function(btns) {
                if(!btns) {
                    btns = this.defaultBtns;
                }
                $('@div.editor-toolbar')
                    .appendChild(this.createSelectionToolbar(btns))
                    .appendChild(this.createInputToolbar())
                    .attachInside('.blog-right-menu');
            },
            createToolbarButtons: function(btns) {
                let buttonList = new DocumentFragment();

                let foundButtons = Object.keys(ToolbarButton).filter( button => {  
                    if(btns == 'all') {
                        return true;
                    }
                    else {
                        for(const value of btns.values()) {
                            if(value == button) {
                                return true;
                            }
                        }
                    }
                });

                foundButtons.forEach( btn=> {
                    if(ToolbarButton[btn].content) {
                        let li = $('@button').addClass(btn)
                            .html(ToolbarButton[btn].content)
                            .click(function() {
                                let selection = window.getSelection();
                                let range = selection.getRangeAt(0).cloneRange();
                                if(ToolbarButton[btn].action) {
                                    if(!ToolbarButton[btn].action.length) {
                                        ToolbarButton[btn].action();
                                    }
                                    if(ToolbarButton[btn].action.length == 1) {
                                        ToolbarButton[btn].action(range);
                                    }
                                }
                            })
                            .appendTo($('@li')).element
                        ;
                        buttonList.appendChild(li);
                    }
                });
                return $("@ul").appendChild(buttonList);
            },
            createSelectionToolbar: function(btns) {
                return $('@div.selectionToolbar').appendChild(this.createToolbarButtons(btns));
            },
            createInputToolbar: function() {
                let inputToolbar =  $('@div.inputToolbar')
                    .appendChild($('@input[type=text placeholder="Enter a link"]'))
                    .appendChild($("@a.save[href=#]").html('<i class="fa fa-check"></i>'))
                    .appendChild($("@a.close[href=#]").html('<i class="fa fa-times"></i>'))
                    .addClass('displaynone');
                return inputToolbar;
            },
        };

        return toolbar;
    })();


    /*  Default Settings  */

    /* init calls */
    let navbar = new Navbar(navtoolbarBtns, navsettingBtns)





    // $('.blog-title').click(function(e) {
    //     /* rgba(30, 30, 30, 0.62) */
    //     let title = $(e.target);
    //     if(title.textContent == 'Add Title' && title.css('color') == "rgba(30, 30, 30, 0.62)") {
    //         title.textContent = '';
    //         title.css('color:#000');
    //     }      
    // });

    // $('.blog-editor').click(function() {
    //     /* Reset Editor to default settings */
    //     if($('.blog-title').textContent == "") {
    //         $('.blog-title').textContent = 'Add Title';
    //         $('.blog-title').css('color:rgba(30, 30, 30, 0.62)');
    //     }
    //     if($('.inputToolbar').css('display') == 'flex') {
    //         $(".inputToolbar").removeClass('displayflex').addClass("displaynone");
    //         $(".selectionToolbar").removeClass('displaynone');
    //     }
    //     // if(window.getSelection().type == 'Range') {
    //     //     $('.editor-toolbar').removeClass('displaynone')
    //     // }
    //     if(window.getSelection().type == 'Caret') {
    //         $('.editor-toolbar').css("visibility:hidden");
    //     }
    // });

    // $('pre code').keydown(function(event) {
    //     let range = Selection.range();

    //     if(event.ctrlKey && event.key == "Enter" && !event.shiftKey|| !event.shiftKey && event.metaKey && event.key == 'Enter') {
    //         console.log("[pressed] ctrl/meta + Enter.");
    //     }
    //     if(event.shiftKey && event.metaKey && event.key== 'Enter') {
    //         console.log("[combinedKey]: insert line before");
    //     }
    //     if(event.key == '(') {
    //         let rightBracket = document.createTextNode(')');
    //         range.insertNode(rightBracket);
    //     }
    //     if(event.key == ')') {
    //         let prevNode = selection.anchorNode;
    //         let nextchar = prevNode.nextSibling.textContent;
    //         if(nextchar == ')') {
    //             prevNode.nextSibling.textContent=""
    //         }
    //     }
    //     if(event.key == '{') {
    //         let rightBracket = document.createTextNode('}');
    //         range.insertNode(rightBracket);
    //     }
    //     if(event.key == '[') {
    //         let rightBracket = document.createTextNode(']');
    //         range.insertNode(rightBracket);
    //     }
    //     if(event.key == 'Backspace') {
    //         event.preventDefault();
    //         console.log("...creating KeyboardEvent");

    //         var kevent = document.createEvent('KeyboardEvent'); 


    //         // if(pairs === '[]' || pairs == '()' || pairs == '{}' || pairs == '""') {
    //         //     console.log("yes: nextSibling = prevsiblign");
    //         //     prevNode.nextSibling.textContent = ''
    //         // }
    //     }
    //     if(event.key == 'Tab') {
    //         event.preventDefault();
    //         let tabStr = document.createTextNode("   ");
    //         range.insertNode(tabStr);
    //         document.getSelection().addRange(range);
    //         let tabrange = window.getSelection().getRangeAt(0).cloneRange();
    //         range.setStart(tabrange.endContainer, tabrange.endOffset);
    //     }
    // });

    /*document.addEventListener('selectionchange',function() {
        let selection = window.getSelection();
        if(selection.type == 'Range' && selection.toString().length > 0) {
            let focusTagName = selection.focusNode.ownerDocument.activeElement.tagName;
            if(focusTagName == 'ARTICLE') {
                let range = selection.getRangeAt(0).cloneRange();
                let rect  = range.getBoundingClientRect();
                let left  = rect.left - $('.editor-toolbar').offsetWidth/2 + rect.width/2;
                let top   = rect.top - $('.editor-toolbar').offsetHeight-8-50;
                $('.editor-toolbar').css(`visibility:visible;left:${left}px;top:${top}px;`)
            }
        }     
    });*/

    // let Editor = (function() {
    //     function editor(editorSelector) {
    //         this.editor = $(editorSelector);
    //     }
    //     editor.prototype = {
    //         clickgroup: function(selectors, callback) {
    //             callback(...selectors);
    //         },
    //         selected: function(listener) {
    //             document.addEventListener('selectionchange', listener);
    //         },
    //         click: function(selector, listener) {
    //             $(selector).click(listener);
    //         },
    //         defaultclick: function(listener) {
    //             this.editor.click(listener)
    //         },
    //         scroll: function(listener) {
    //             this.editor.addEventListener('scroll', listener);
    //         },
    //     };
    //     return editor;
    // })();

})();