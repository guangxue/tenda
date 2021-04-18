let Blog = (function() {

    let Element = (function() {


        let constructor = function(abbr) {
            let classNames = abbr.match(/(?<=\.)[0-9a-zA-Z-]+/g);
            let idvalue = abbr.match(/(?<=#)[a-zA-Z-]+/) && abbr.match(/(?<=#)[a-zA-Z-]+/)[0];
            let attrStart = abbr.indexOf('[');
            let attrEnd = abbr.lastIndexOf(']');
            let tagName = abbr.split(/[\.|#|\[]/)[0];

            this.element = document.createElement(tagName);
            if(classNames) {
                classNames.forEach(cls=>{
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
                attrs = attrs.replaceAll(/[0-9a-zA-z-#_]+(?=:)/g, '"$&"');
                if(attrs.slice(-2, -1) === ',') {
                    attrs = attrs.slice(0, -2)+'}'
                }
                let attrObj = JSON.parse(attrs)
                for(const key in attrObj) {
                    this.element.setAttribute(key, attrObj[key])
                }
            }
        };

        constructor.prototype = {
            appendChild: function(childElement) {
                if(childElement.hasOwnProperty('element')) {
                    this.element.appendChild(childElement.element);
                }
                else {
                    this.element.appendChild(childElement);
                }
                return this;
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
            innerHTML: function(htmlstr) {
                this.element.innerHTML = htmlstr;
                return this;
            },
            attachAfter: function(liveElement) {
                $(liveElement).insertAdjacentElement('afterend', this.element);
            },
        };

        return constructor;
    })();

    let $ = function(selector) {
        try
        {
            let element = document.querySelector(selector);
            element.on = function(event, listener) {
               element.addEventListener(event, listener);
            }
            element.click = function(listener) {
                element.addEventListener('click', listener);
            }
            element.keyup = function(listener) {
                element.addEventListener('keyup', listener);
            }
            element.keydown = function(listener) {
                element.addEventListener('keydown', listener);
            }
            element.mouseup = function(listener) {
                element.addEventListener('mouseup', listener);
            }
            element.css = function(...styles) {
                if(styles.length == 1) {

                }
                if(styles.length == 2) {

                }
            }
            return element;
        } catch(e) {
            return new createElement(selector)
        }
    };

    let ToolbarButton = {
            'bold': {
                CommandName: 'bold',
                ValueArgs: '',
                content: '<i class="fa fa-bold" aria-hidden="true"></i>'
            },
            'action': {
                'save': '<i class="fa fa-check"></i>',
                'close': '<i class="fa fa-times"></i>',
            },
            'italic': {
                CommandName: 'italic',
                ValueArgs: '',
                content: '<i class="fa fa-italic" aria-hidden="true"></i>'
            },
            'underline': {
                CommandName: 'underline',
                ValueArgs: '',
                content: '<i class="fa fa-underline" aria-hidden="true"></i>'
            },
            'strikethrough': {
                CommandName: 'strikeThrough',
                ValueArgs: '',
                content: '<i class="fa fa-strikethrough" aria-hidden="true"></i>'
            },
            'superscript': {
                CommandName: 'superscript',
                ValueArgs: '',
                content: '<i class="fa fa-superscript" aria-hidden="true"></i>'
            },
            'subscript': {
                CommandName: 'subscript',
                ValueArgs: '',
                content: '<i class="fa fa-subscript" aria-hidden="true"></i>'
            },
            'link': {
                CommandName: 'createLink',
                ValueArgs: 'http',
                placeholder: 'Enter a link',
                content: '<i class="fa fa-link" aria-hidden="true"></i>'
            },
            'image': {
                CommandName: 'insertImage',
                ValueArgs: 'image',
                placeholder: 'Enter a image link',
                content: '<i class="fa fa-picture-o" aria-hidden="true"></i>'
            },
            'html': {
                CommandName: 'code',
                content: '<i class="fa fa-code" aria-hidden="true"></i>'
            },
            'table': {
                CommandName: 'table',
                content: '<i class="fa fa-table" aria-hidden="true"></i>',
            },
            'orderedlist': {
                CommandName: 'insertOrderedList',
                content: '<i class="fa fa-list-ol" aria-hidden="true"></i>'
            },
            'unorderedlist': {
                CommandName: 'insertUnorderedList',
                content: '<i class="fa fa-list-ul" aria-hidden="true"></i>'
            },
            'indent': {
                CommandName: 'indent',
                content: '<i class="fa fa-indent" aria-hidden="true"></i>'
            },
            'outdent': {
                CommandName: 'outdent',
                content: '<i class="fa fa-outdent" aria-hidden="true"></i>'
            },
            'justifyCenter': {
                CommandName: 'justifyCenter',
                content: '<i class="fa fa-align-center" aria-hidden="true"></i>'
            },
            'justifyFull': {
                CommandName: 'justifyFull',
                content: '<i class="fa fa-align-justify" aria-hidden="true"></i>'
            },
            'justifyLeft': {
                CommandName: 'justifyLeft',
                content: '<i class="fa fa-align-left" aria-hidden="true"></i>'
            },
            'justifyRight': {
                CommandName: 'justifyRight',
                content: '<i class="fa fa-align-right" aria-hidden="true"></i>'
            },
            // Known inline elements that are not removed, or not removed consistantly across browsers:
            // <span>, <label>, <br>
            'removeFormat': {
                CommandName: 'removeFormat',
                content: '<i class="fa fa-eraser" aria-hidden="true"></i>'
            },
            'quote': {
                CommandName: 'quote',
                content: '<i class="fa fa-quote-right" aria-hidden="true"></i>'
            },
            'code': {
                CommandName: 'pre',
                content: '<i class="fa fa-code fa-lg" aria-hidden="true"></i>'
            },
            'h1': {
                CommandName: 'heading',
                ValueArgs: 'H1',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>1</sup>'
            },
            'h2': {
                CommandName: 'heading',
                ValueArgs: 'H2',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>2</sup>'
            },
            'h3': {
                CommandName: 'heading',
                ValueArgs: 'H3',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>3</sup>'
            },
            'h4': {
                CommandName: 'heading',
                ValueArgs: 'H4',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>4</sup>'
            },
            'h5': {
                CommandName: 'heading',
                ValueArgs: 'H5',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>5</sup>'
            },
            'h6': {
                CommandName: 'heading',
                ValueArgs: 'H6',
                content: '<i class="fa fa-header" aria-hidden="true"><sup>6</sup>'
            }
    };

    let Toolbar = (function(buttons) {

        let toolbar = function(btns) {
            return this.init(btns)
        };
        toolbar.prototype = {
            init: function(btns) {
                return this.createToolbar(btns);
            },
            createToolbar: function(btns) {
                if(!btns) {
                    btns = ['bold', 'italic', 'underline', 'link', 'h2', 'h3', 'quote'];
                }
                $('div.editor-toolbar')
                    .appendChild($('div.toolbarWrapper').appendChild(this.createToolbarButtons(btns)))
                    .appendChild($('div.inputWrapper').appendChild(this.createToolbarInputs()))
                    .attachAfter('.blog-container');
            },
            createToolbarButtons: function(btns) {
                $('ul').append
                let ul = document.createElement('ul');
                let ButtonListFragment = new DocumentFragment();

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
                foundButtons.forEach((btn) => {
                    let li = document.createElement('li');
                    
                    if(ToolbarButton[btn].content) {
                        let button = document.createElement('button');
                        button.className = btn;
                        button.innerHTML = ToolbarButton[btn].content;
                        li.appendChild(button);
                    }
                    
                    ButtonListFragment.appendChild(li)
                });
                ul.appendChild(ButtonListFragment);
                return ul;
            },
            createToolbarInputs: function() {
                let inputWrapper =  $('div.inputWrapper')
                    .appendChild($('input[type=text placeholder="Enter a link"]'))
                    .appendChild($("a.save[href=#]").innerHTML(ToolbarButton.action['close']))
                    .appendChild($("a.close[href=#]").innerHTML(ToolbarButton.action['close']));
                return inputWrapper;
            },
        };
        return toolbar;
    })();

    let Selection = {
        clientRect: function() {
            if(window.getSelection().type === 'Range') {
                let selection = window.getSelection();
                let range = selection.getRangeAt(0).cloneRange();
                let DOMRect = range.getBoundingClientRect();
                if(DOMRect.width) {
                    return DOMRect;
                }
            }
            else {
                return;
            }
        },
    };

    let Editor = (function() {
        function editor(editorSelector) {
            this.editor = $(editorSelector);
        }
        editor.prototype = {
            click: function(selectors, callback) {
                callback(...selectors);
            },
        };
        return editor;
    })();

    /*  Default Settings  */
    let toolbar = new Toolbar();
    let editor  = new Editor('.editor');
    editor.click([$('.blog-title'), $('.editor')], function(title, editor) {
        title.click(function() {
            title.innerHTML = ""
        })
        editor.click(function() {
            if(title.textContent == "") {
                title.textContent = 'Add Title';
            }
        })
    })
    return {
        toolbar: Toolbar,
    }
})(); 






