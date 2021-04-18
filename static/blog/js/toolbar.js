import {
	FAButton,
	ElementWrapper as $,
	Selection,
} from './core.js';

function Toolbar() {
	this.defaultBtns = ['bold', 'italic', 'underline', 'link', 'h2', 'h3', 'quote']
	return this.init();
}
Toolbar.prototype = FAButton;
Toolbar.prototype.init = function() {
	let selectionToolbar = FAButton.build(this.defaultBtns).wrapperElement('div').wrapperClass('selectToolbar');
	$('.blog-container').addChild($('@div.toolbar-container'));
	$('.toolbar-container').addChild(selectionToolbar.wrapper).addChild(this.creatInputToolbar())
};
Toolbar.prototype.rebuild = function(buttons) {
	let newSelectBar = FAButton.New().build(buttons).wrapperElement('div').wrapperClass('selectToolbar');
	$('.toolbar-container .selectToolbar').replaceWith(newSelectBar.wrapper)
};
Toolbar.prototype.creatInputToolbar = function() {
	let inputToolbarWrapper =  $('@div.inputToolbar')
	        .addChild($('@input[type=text placeholder="Enter a link"]'))
	        .addChild($("@a.save[href=#]").html(FAButton.getButton('check')))
	        .addChild($("@a.close[href=#]").html(FAButton.getButton('times')))
	        .addClass('displaynone');
	this.inputToolbar = inputToolbarWrapper.element;
	return inputToolbarWrapper;
};
Toolbar.prototype.inputButton = function(elem) {
	let el = this.inputToolbar.querySelector(elem)
	return $(el);
};

export default Toolbar;

