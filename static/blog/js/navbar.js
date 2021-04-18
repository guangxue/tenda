import {
	ElementWrapper as $,
	WPButton,
} from './core.js';

function Navbar() {
	return this.init();
};

Navbar.prototype = WPButton;
Navbar.prototype.init = function() {
	let navtoolbar = $('@div.navtoolbar');
	let navtoolbarButtons = WPButton.buildwith(['plus', 'undo', 'redo'], 'div');
	navtoolbar.addChild(navtoolbarButtons.allButtons);
	let middleflex = $('@div.middleflex');
	let navsetting = $('@div.navsetting');
	let previewButton = $('@div').addChild($('@button.preview').text('Preview'))
	let publishButton = $('@div').addChild($('@button.publish').text('Publish'))
	let settingButton = $("@div").addChild($('@button.setting').addChild($('@span').html(WPButton.getButton('setting'))))
	navsetting.addChild(previewButton).addChild(publishButton).addChild(settingButton);
	$('header nav').addChild(navtoolbar).addChild(middleflex).addChild(navsetting);
	$('header nav button.plus').addClass('navbutton-plus');
	$('.blog-right-menu').addClass('displaynone');
};

let navbar = new Navbar();


export default Navbar;