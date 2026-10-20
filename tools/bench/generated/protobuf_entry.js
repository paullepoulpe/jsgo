/** Actual protoc-generated descriptor and well-known-type round trips. */
const descriptor = require('google-protobuf/google/protobuf/descriptor_pb.js');
const timestamp = require('google-protobuf/google/protobuf/timestamp_pb.js');
const struct = require('google-protobuf/google/protobuf/struct_pb.js');

/** @return {string} */
function exerciseProtobuf() {
  const field = new descriptor.FieldDescriptorProto();
  field.setName('total');
  field.setNumber(1);
  field.setType(descriptor.FieldDescriptorProto.Type.TYPE_INT64);
  const message = new descriptor.DescriptorProto();
  message.setName('Invoice');
  message.addField(field);
  const file = new descriptor.FileDescriptorProto();
  file.setName('enterprise.proto');
  file.setPackage('jsgo.fixture');
  file.addMessageType(message);
  const files = new descriptor.FileDescriptorSet();
  files.addFile(file);
  const decoded = descriptor.FileDescriptorSet.deserializeBinary(files.serializeBinary());
  const decodedFile = decoded.getFileList()[0];
  const decodedMessage = decodedFile.getMessageTypeList()[0];
  const decodedField = decodedMessage.getFieldList()[0];
  const time = new timestamp.Timestamp();
  time.setSeconds(1700000000);
  time.setNanos(123456789);
  const decodedTime = timestamp.Timestamp.deserializeBinary(time.serializeBinary());
  const object = struct.Struct.fromJavaScript({'name': 'Ada', 'active': true, 'values': [3, 'x']});
  const decodedObject = struct.Struct.deserializeBinary(object.serializeBinary()).toJavaScript();
  if (decodedField.getNumber() !== 1 ||
      decodedField.getType() !== descriptor.FieldDescriptorProto.Type.TYPE_INT64 ||
      decodedObject['active'] !== true || decodedObject['values'][1] !== 'x') {
    throw new Error('protoc-generated serialization semantics failed');
  }
  return decodedFile.getName() + '|' + decodedMessage.getName() + '|' +
      decodedField.getName() + '|' + decodedTime.getSeconds() + '|' +
      decodedTime.getNanos() + '|' + decodedObject['name'];
}

globalThis['jsgoEnterpriseProtobuf'] = exerciseProtobuf;
