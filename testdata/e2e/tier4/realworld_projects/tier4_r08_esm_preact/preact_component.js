function h(tag, props, ...children) { return { tag, props, children }; }
const Component = (props) => h('div', { class: 'btn' }, props.label);
window['render'] = Component({ label: 'Click' });
