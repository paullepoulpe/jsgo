/**
 * @fileoverview Benchmark sample app (tools/bench/README.md): the notes UI, a
 * tab per note, a rich-text editor with toolbar and plugins (including the
 * Closure Library demo plugins), a due-date picker and a sync button.
 */
goog.module('benchapp.NotesView');

const BasicTextFormatter = goog.require('goog.editor.plugins.BasicTextFormatter');
const Button = goog.require('goog.ui.Button');
const Command = goog.require('goog.editor.Command');
const Component = goog.require('goog.ui.Component');
const DatePicker = goog.require('goog.ui.DatePicker');
const DateTimeFormat = goog.require('goog.i18n.DateTimeFormat');
const DefaultToolbar = goog.require('goog.ui.editor.DefaultToolbar');
const Delay = goog.require('goog.async.Delay');
const EnterHandler = goog.require('goog.editor.plugins.EnterHandler');
const Field = goog.require('goog.editor.Field');
const HeaderFormatter = goog.require('goog.editor.plugins.HeaderFormatter');
const HelloWorld = goog.require('goog.demos.editor.HelloWorld');
const HelloWorldDialogPlugin = goog.require('goog.demos.editor.HelloWorldDialogPlugin');
const LinkBubble = goog.require('goog.editor.plugins.LinkBubble');
const LinkDialogPlugin = goog.require('goog.editor.plugins.LinkDialogPlugin');
const ListTabHandler = goog.require('goog.editor.plugins.ListTabHandler');
const NoteStore = goog.require('benchapp.NoteStore');
const RemoveFormatting = goog.require('goog.editor.plugins.RemoveFormatting');
const SafeHtml = goog.require('goog.html.SafeHtml');
const SpacesTabHandler = goog.require('goog.editor.plugins.SpacesTabHandler');
const Tab = goog.require('goog.ui.Tab');
const TabBar = goog.require('goog.ui.TabBar');
const TagName = goog.require('goog.dom.TagName');
const ToolbarController = goog.require('goog.ui.editor.ToolbarController');
const UndoRedo = goog.require('goog.editor.plugins.UndoRedo');
const XhrIo = goog.require('goog.net.XhrIo');
const dom = goog.require('goog.dom');
const events = goog.require('goog.events');

/** @const {!Array<string>} */
const TOOLBAR_BUTTONS = [
  Command.BOLD, Command.ITALIC, Command.UNDERLINE, Command.FONT_COLOR,
  Command.FONT_FACE, Command.LINK, Command.UNDO, Command.REDO,
  Command.UNORDERED_LIST, Command.ORDERED_LIST, Command.INDENT,
  Command.OUTDENT, Command.REMOVE_FORMAT,
];

/** The notes UI over a NoteStore. */
class NotesView extends Component {
  /**
   * @param {!NoteStore} store
   * @param {string} syncUrl Where "Sync" posts the notes.
   */
  constructor(store, syncUrl) {
    super();

    /** @private @const {!NoteStore} */
    this.store_ = store;

    /** @private @const {string} */
    this.syncUrl_ = syncUrl;

    /** @private @const {!TabBar} */
    this.tabs_ = new TabBar();

    /** @private @const {!DatePicker} */
    this.datePicker_ = new DatePicker();

    /** @private @const {!Button} */
    this.syncButton_ = new Button('Sync');

    /** @private @const {!DateTimeFormat} */
    this.dateFormat_ = new DateTimeFormat(DateTimeFormat.Format.MEDIUM_DATE);

    /** @private {?Field} */
    this.field_ = null;

    /** @private {?ToolbarController} */
    this.toolbarController_ = null;

    /** @private {?string} */
    this.currentId_ = null;

    /** @private @const {!Delay} Saves the editor contents after typing stops. */
    this.saveDelay_ = new Delay(() => this.saveCurrent_(), 500);

    this.addChild(this.tabs_);
    this.addChild(this.datePicker_);
    this.addChild(this.syncButton_);
  }

  /** @override */
  createDom() {
    const d = this.getDomHelper();
    const element = d.createDom(
        TagName.DIV, 'notes',
        d.createDom(TagName.DIV, 'notes-tabs'),
        d.createDom(TagName.DIV, {'id': 'notes-toolbar'}),
        d.createDom(TagName.DIV, {'id': 'notes-editor'}),
        d.createDom(TagName.DIV, 'notes-side'),
        d.createDom(TagName.SPAN, 'notes-status'));
    this.setElementInternal(element);
  }

  /** @override */
  enterDocument() {
    super.enterDocument();
    const el = this.getElement();
    this.tabs_.render(dom.getElementByClass('notes-tabs', el));
    this.datePicker_.render(dom.getElementByClass('notes-side', el));
    this.syncButton_.render(dom.getElementByClass('notes-side', el));

    this.field_ = this.createEditor_();
    this.getHandler()
        .listen(this.store_, NoteStore.CHANGE, this.refreshTabs_)
        .listen(this.tabs_, Component.EventType.SELECT, this.onSelectTab_)
        .listen(this.datePicker_, DatePicker.Events.CHANGE, this.onDateChange_)
        .listen(this.syncButton_, Component.EventType.ACTION, this.sync_)
        .listen(this.field_, Field.EventType.DELAYEDCHANGE, () => this.saveDelay_.start());
    this.field_.makeEditable();
    this.refreshTabs_();
  }

  /**
   * @return {!Field}
   * @private
   */
  createEditor_() {
    const field = new Field('notes-editor');
    field.registerPlugin(new BasicTextFormatter());
    field.registerPlugin(new RemoveFormatting());
    field.registerPlugin(new UndoRedo());
    field.registerPlugin(new ListTabHandler());
    field.registerPlugin(new SpacesTabHandler());
    field.registerPlugin(new EnterHandler());
    field.registerPlugin(new HeaderFormatter());
    field.registerPlugin(new LinkDialogPlugin());
    field.registerPlugin(new LinkBubble());
    field.registerPlugin(new HelloWorld());
    field.registerPlugin(new HelloWorldDialogPlugin());
    const toolbar = DefaultToolbar.makeToolbar(
        TOOLBAR_BUTTONS, dom.getRequiredElement('notes-toolbar'));
    this.toolbarController_ = new ToolbarController(field, toolbar);
    return field;
  }

  /** @private */
  refreshTabs_() {
    this.tabs_.removeChildren(true);
    for (const note of this.store_.getAll()) {
      const tab = new Tab(note.title);
      tab.setId(note.id);
      this.tabs_.addChild(tab, true);
      if (note.id == this.currentId_) {
        this.tabs_.setSelectedTab(tab);
      }
    }
    this.setStatus_(this.store_.getAll().length + ' notes');
  }

  /**
   * @param {!events.Event} e
   * @private
   */
  onSelectTab_(e) {
    const tab = /** @type {!Tab} */ (e.target);
    this.saveCurrent_();
    this.currentId_ = tab.getId();
    const note = this.store_.get(this.currentId_);
    if (note && this.field_) {
      this.field_.setSafeHtml(false, SafeHtml.htmlEscape(note.html));
      this.datePicker_.setDate(note.due == null ? null : new Date(note.due));
    }
  }

  /** @private */
  onDateChange_() {
    const date = this.datePicker_.getDate();
    const id = this.currentId_;
    if (id == null) {
      return;
    }
    this.store_.update(id, (note) => {
      note.due = date ? date.valueOf() : null;
    });
    if (date) {
      this.setStatus_('Due ' + this.dateFormat_.format(date));
    }
  }

  /** @private */
  saveCurrent_() {
    const id = this.currentId_;
    const field = this.field_;
    if (id == null || !field) {
      return;
    }
    const html = field.getCleanContents();
    this.store_.update(id, (note) => {
      note.html = html;
    });
  }

  /** @private */
  sync_() {
    this.saveCurrent_();
    this.setStatus_('Syncing...');
    XhrIo.send(
        this.syncUrl_, (e) => {
          const xhr = /** @type {!XhrIo} */ (e.target);
          this.setStatus_(xhr.isSuccess() ? 'Synced' : 'Sync failed: ' + xhr.getStatus());
        },
        'POST', JSON.stringify(this.store_.getAll()),
        {'Content-Type': 'application/json'});
  }

  /**
   * @param {string} text
   * @private
   */
  setStatus_(text) {
    dom.setTextContent(dom.getElementByClass('notes-status', this.getElement()), text);
  }

  /** @override */
  disposeInternal() {
    this.saveDelay_.dispose();
    if (this.toolbarController_) {
      this.toolbarController_.dispose();
    }
    if (this.field_) {
      this.field_.dispose();
    }
    super.disposeInternal();
  }
}

exports = NotesView;
