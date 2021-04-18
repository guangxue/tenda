import toolbar from './toolbar.js'
import leftsidebar from './sidebar.js'
import { ElementWrapper as $, Selection } from './core.js'

function Editor(editorElement) {
	this.editor = editorElement;
}

Editor.prototype = {
	click: function(listener) {
		$(this.editor).click(listener)
	},
	selected: function(callback) {
		document.addEventListener('selectionchange', function() {
            let sel = window.getSelection();
            if(sel.anchorNode.ownerDocument.activeElement.tagName == 'ARTICLE') {
                callback();
            }
            else {
                return false;
            }
        });
	},
};

export default Editor;






