# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: orderer/kafka.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='orderer/kafka.proto',
  package='orderer',
  syntax='proto3',
  serialized_pb=_b('\n\x13orderer/kafka.proto\x12\x07orderer\"\xaf\x01\n\x0cKafkaMessage\x12/\n\x07regular\x18\x01 \x01(\x0b\x32\x1c.orderer.KafkaMessageRegularH\x00\x12\x35\n\x0btime_to_cut\x18\x02 \x01(\x0b\x32\x1e.orderer.KafkaMessageTimeToCutH\x00\x12/\n\x07\x63onnect\x18\x03 \x01(\x0b\x32\x1c.orderer.KafkaMessageConnectH\x00\x42\x06\n\x04Type\"&\n\x13KafkaMessageRegular\x12\x0f\n\x07payload\x18\x01 \x01(\x0c\"-\n\x15KafkaMessageTimeToCut\x12\x14\n\x0c\x62lock_number\x18\x01 \x01(\x04\"&\n\x13KafkaMessageConnect\x12\x0f\n\x07payload\x18\x01 \x01(\x0c\".\n\rKafkaMetadata\x12\x1d\n\x15last_offset_persisted\x18\x01 \x01(\x03\x42U\n%org.hyperledger.fabric.protos.ordererZ,github.com/hyperledger/fabric/protos/ordererb\x06proto3')
)
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_KAFKAMESSAGE = _descriptor.Descriptor(
  name='KafkaMessage',
  full_name='orderer.KafkaMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='regular', full_name='orderer.KafkaMessage.regular', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='time_to_cut', full_name='orderer.KafkaMessage.time_to_cut', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='connect', full_name='orderer.KafkaMessage.connect', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='Type', full_name='orderer.KafkaMessage.Type',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=33,
  serialized_end=208,
)


_KAFKAMESSAGEREGULAR = _descriptor.Descriptor(
  name='KafkaMessageRegular',
  full_name='orderer.KafkaMessageRegular',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='payload', full_name='orderer.KafkaMessageRegular.payload', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=210,
  serialized_end=248,
)


_KAFKAMESSAGETIMETOCUT = _descriptor.Descriptor(
  name='KafkaMessageTimeToCut',
  full_name='orderer.KafkaMessageTimeToCut',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='block_number', full_name='orderer.KafkaMessageTimeToCut.block_number', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=250,
  serialized_end=295,
)


_KAFKAMESSAGECONNECT = _descriptor.Descriptor(
  name='KafkaMessageConnect',
  full_name='orderer.KafkaMessageConnect',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='payload', full_name='orderer.KafkaMessageConnect.payload', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=297,
  serialized_end=335,
)


_KAFKAMETADATA = _descriptor.Descriptor(
  name='KafkaMetadata',
  full_name='orderer.KafkaMetadata',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='last_offset_persisted', full_name='orderer.KafkaMetadata.last_offset_persisted', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=337,
  serialized_end=383,
)

_KAFKAMESSAGE.fields_by_name['regular'].message_type = _KAFKAMESSAGEREGULAR
_KAFKAMESSAGE.fields_by_name['time_to_cut'].message_type = _KAFKAMESSAGETIMETOCUT
_KAFKAMESSAGE.fields_by_name['connect'].message_type = _KAFKAMESSAGECONNECT
_KAFKAMESSAGE.oneofs_by_name['Type'].fields.append(
  _KAFKAMESSAGE.fields_by_name['regular'])
_KAFKAMESSAGE.fields_by_name['regular'].containing_oneof = _KAFKAMESSAGE.oneofs_by_name['Type']
_KAFKAMESSAGE.oneofs_by_name['Type'].fields.append(
  _KAFKAMESSAGE.fields_by_name['time_to_cut'])
_KAFKAMESSAGE.fields_by_name['time_to_cut'].containing_oneof = _KAFKAMESSAGE.oneofs_by_name['Type']
_KAFKAMESSAGE.oneofs_by_name['Type'].fields.append(
  _KAFKAMESSAGE.fields_by_name['connect'])
_KAFKAMESSAGE.fields_by_name['connect'].containing_oneof = _KAFKAMESSAGE.oneofs_by_name['Type']
DESCRIPTOR.message_types_by_name['KafkaMessage'] = _KAFKAMESSAGE
DESCRIPTOR.message_types_by_name['KafkaMessageRegular'] = _KAFKAMESSAGEREGULAR
DESCRIPTOR.message_types_by_name['KafkaMessageTimeToCut'] = _KAFKAMESSAGETIMETOCUT
DESCRIPTOR.message_types_by_name['KafkaMessageConnect'] = _KAFKAMESSAGECONNECT
DESCRIPTOR.message_types_by_name['KafkaMetadata'] = _KAFKAMETADATA

KafkaMessage = _reflection.GeneratedProtocolMessageType('KafkaMessage', (_message.Message,), dict(
  DESCRIPTOR = _KAFKAMESSAGE,
  __module__ = 'orderer.kafka_pb2'
  # @@protoc_insertion_point(class_scope:orderer.KafkaMessage)
  ))
_sym_db.RegisterMessage(KafkaMessage)

KafkaMessageRegular = _reflection.GeneratedProtocolMessageType('KafkaMessageRegular', (_message.Message,), dict(
  DESCRIPTOR = _KAFKAMESSAGEREGULAR,
  __module__ = 'orderer.kafka_pb2'
  # @@protoc_insertion_point(class_scope:orderer.KafkaMessageRegular)
  ))
_sym_db.RegisterMessage(KafkaMessageRegular)

KafkaMessageTimeToCut = _reflection.GeneratedProtocolMessageType('KafkaMessageTimeToCut', (_message.Message,), dict(
  DESCRIPTOR = _KAFKAMESSAGETIMETOCUT,
  __module__ = 'orderer.kafka_pb2'
  # @@protoc_insertion_point(class_scope:orderer.KafkaMessageTimeToCut)
  ))
_sym_db.RegisterMessage(KafkaMessageTimeToCut)

KafkaMessageConnect = _reflection.GeneratedProtocolMessageType('KafkaMessageConnect', (_message.Message,), dict(
  DESCRIPTOR = _KAFKAMESSAGECONNECT,
  __module__ = 'orderer.kafka_pb2'
  # @@protoc_insertion_point(class_scope:orderer.KafkaMessageConnect)
  ))
_sym_db.RegisterMessage(KafkaMessageConnect)

KafkaMetadata = _reflection.GeneratedProtocolMessageType('KafkaMetadata', (_message.Message,), dict(
  DESCRIPTOR = _KAFKAMETADATA,
  __module__ = 'orderer.kafka_pb2'
  # @@protoc_insertion_point(class_scope:orderer.KafkaMetadata)
  ))
_sym_db.RegisterMessage(KafkaMetadata)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('\n%org.hyperledger.fabric.protos.ordererZ,github.com/hyperledger/fabric/protos/orderer'))
try:
  # THESE ELEMENTS WILL BE DEPRECATED.
  # Please use the generated *_pb2_grpc.py files instead.
  import grpc
  from grpc.framework.common import cardinality
  from grpc.framework.interfaces.face import utilities as face_utilities
  from grpc.beta import implementations as beta_implementations
  from grpc.beta import interfaces as beta_interfaces
except ImportError:
  pass
# @@protoc_insertion_point(module_scope)
