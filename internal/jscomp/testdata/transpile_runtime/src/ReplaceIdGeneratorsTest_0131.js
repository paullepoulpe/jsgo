/** @idGenerator {stable} */ function tagId(strings, ...args) {return strings[0];};
function Foo() { tagId`foo${1}bar${2}`; }
