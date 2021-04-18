import Toolbar from './toolbar.js';
import Sidebar from './sidebar.js';
import Editor from './editor.js';
import Navbar from './navbar.js';
import { ElementWrapper as $, Selection } from './core.js';


/* LeftSideBar init */
let leftsidebar = new Sidebar('left');
leftsidebar
	.addHeader('HEADER')
	.addTitle('text')
	.addButtons(['image', 'code', 'table', 'earser'])
	.addTitle('design')

leftsidebar.button('button.code').click( function() {
	console.log("...adding Code Editor");
	let codeWrapper = $('@pre').css('border: 1px solid black;padding: 1em 1em;')
        				.addChild($('@code').css('padding:5px;'));
    let codeEditor = $('@div.codeEditor').addChild(codeWrapper);
    let parentNode = Selection.wrapperNode();
    $(parentNode).addSibling(codeEditor)
});



/* Toolbar Init */
let toolbar = new Toolbar();
toolbar.button('button.link').click(function() {
	Selection.removeAllRanges();
	Selection.addRange(Selection.range)
	$('div.selectToolbar').addClass('displaynone');
	$('div.inputToolbar').removeClass('displaynone').addClass('displayflex');
	
});
toolbar.inputButton('a.close').click(function() {
	$('div.inputToolbar').removeClass('displayflex').addClass('displaynone');
	$('div.selectToolbar').removeClass('displaynone');
});


/* Editor init */
let editor = new Editor('article.blog-editor');
editor.selected(function() {
	if(!Selection.position()) { return }
	let pos  = Selection.position();
	let paddingTop = 68;
	let halfToolbarWidth = $('.toolbar-container').offsetWidth/2;
	let halfSelectionWidth = pos.width/2;
	let toolbarHeight = $('.toolbar-container').offsetHeight;
	let left = pos.left - halfToolbarWidth + halfSelectionWidth;
	let top  = pos.top - toolbarHeight - paddingTop;
	// console.log("* rectTop:", pos.top);
	// console.log("* toolbarHeight:", toolbarHeight);
	// console.log("* paddingTop:", paddingTop);
	// console.log("=> top now:", top);
	// console.log("isRange():", Selection.isRange());
	if(Selection.isRange()) {
		$(".toolbar-container").css(`visibility:visible;left:${left}px;top:${top}px;`)
	}
});

editor.click(function() {
	if(!Selection.isRange()) {
		$(".toolbar-container").css("visibility:hidden")
	}
});


$('button.plus').click(function() {
	$('.blog-left-menu').toggleClass('displaynone');
});

$('button.setting').click(function() {
	if(!$('.blog-left-menu').hasClass('displaynone')) {
		$('.blog-left-menu').addClass('displaynone');
		$('.blog-right-menu').removeClass('displaynone');
	}
	else {
		$('.blog-right-menu').toggleClass('displaynone');
	}
});

export {
	toolbar,
	leftsidebar,
	editor,
}