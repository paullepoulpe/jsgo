/** Benchmark-owned semantic harness over genuine upstream generated messages. */
goog.provide('enterprise.protobufClosure');

goog.require('proto.google.protobuf.DescriptorProto');
goog.require('proto.google.protobuf.FieldDescriptorProto');
goog.require('proto.google.protobuf.FileDescriptorProto');
goog.require('proto.google.protobuf.FileDescriptorSet');
goog.require('proto.google.protobuf.ListValue');
goog.require('proto.google.protobuf.Struct');
goog.require('proto.google.protobuf.Timestamp');
goog.require('proto.google.protobuf.Value');

/** @return {string} */
enterprise.protobufClosure.exercise = function() {
  const field = new proto.google.protobuf.FieldDescriptorProto();
  field.setName('total');
  field.setNumber(1);
  field.setType(proto.google.protobuf.FieldDescriptorProto.Type.TYPE_INT64);
  const message = new proto.google.protobuf.DescriptorProto();
  message.setName('Invoice');
  message.addField(field);
  const file = new proto.google.protobuf.FileDescriptorProto();
  file.setName('enterprise.proto');
  file.setPackage('jsgo.fixture');
  file.addMessageType(message);
  const files = new proto.google.protobuf.FileDescriptorSet();
  files.addFile(file);
  const decoded = proto.google.protobuf.FileDescriptorSet.deserializeBinary(files.serializeBinary());
  const decodedFile = decoded.getFileList()[0];
  const decodedMessage = decodedFile.getMessageTypeList()[0];
  const decodedField = decodedMessage.getFieldList()[0];
  const time = new proto.google.protobuf.Timestamp();
  time.setSeconds(1700000000);
  time.setNanos(123456789);
  const decodedTime = proto.google.protobuf.Timestamp.deserializeBinary(time.serializeBinary());
  const object = new proto.google.protobuf.Struct();
  const name = new proto.google.protobuf.Value();
  name.setStringValue('Ada');
  const active = new proto.google.protobuf.Value();
  active.setBoolValue(true);
  const item = new proto.google.protobuf.Value();
  item.setStringValue('x');
  const list = new proto.google.protobuf.ListValue();
  list.addValues(item);
  const values = new proto.google.protobuf.Value();
  values.setListValue(list);
  object.getFieldsMap().set('name', name);
  object.getFieldsMap().set('active', active);
  object.getFieldsMap().set('values', values);
  const decodedObject = proto.google.protobuf.Struct.deserializeBinary(object.serializeBinary());
  if (decodedField.getNumber() !== 1 ||
      decodedField.getType() !== proto.google.protobuf.FieldDescriptorProto.Type.TYPE_INT64 ||
      !decodedObject.getFieldsMap().get('active').getBoolValue() ||
      decodedObject.getFieldsMap().get('values').getListValue().getValuesList()[0].getStringValue() !== 'x') {
    throw new Error('Closure-style generated protobuf binary round trip failed');
  }
  return decodedFile.getName() + '|' + decodedMessage.getName() + '|' +
      decodedField.getName() + '|' + decodedTime.getSeconds() + '|' +
      decodedTime.getNanos() + '|' + decodedObject.getFieldsMap().get('name').getStringValue();
};

globalThis['jsgoEnterpriseProtobufClosure'] = enterprise.protobufClosure.exercise;
