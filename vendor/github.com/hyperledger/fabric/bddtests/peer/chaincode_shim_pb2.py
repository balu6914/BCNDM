# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: peer/chaincode_shim.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from peer import chaincode_event_pb2 as peer_dot_chaincode__event__pb2
from peer import proposal_pb2 as peer_dot_proposal__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='peer/chaincode_shim.proto',
  package='protos',
  syntax='proto3',
  serialized_pb=_b('\n\x19peer/chaincode_shim.proto\x12\x06protos\x1a\x1apeer/chaincode_event.proto\x1a\x13peer/proposal.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xac\x04\n\x10\x43haincodeMessage\x12+\n\x04type\x18\x01 \x01(\x0e\x32\x1d.protos.ChaincodeMessage.Type\x12-\n\ttimestamp\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0f\n\x07payload\x18\x03 \x01(\x0c\x12\x0c\n\x04txid\x18\x04 \x01(\t\x12\"\n\x08proposal\x18\x05 \x01(\x0b\x32\x10.protos.Proposal\x12/\n\x0f\x63haincode_event\x18\x06 \x01(\x0b\x32\x16.protos.ChaincodeEvent\"\xc7\x02\n\x04Type\x12\r\n\tUNDEFINED\x10\x00\x12\x0c\n\x08REGISTER\x10\x01\x12\x0e\n\nREGISTERED\x10\x02\x12\x08\n\x04INIT\x10\x03\x12\t\n\x05READY\x10\x04\x12\x0f\n\x0bTRANSACTION\x10\x05\x12\r\n\tCOMPLETED\x10\x06\x12\t\n\x05\x45RROR\x10\x07\x12\r\n\tGET_STATE\x10\x08\x12\r\n\tPUT_STATE\x10\t\x12\r\n\tDEL_STATE\x10\n\x12\x14\n\x10INVOKE_CHAINCODE\x10\x0b\x12\x0c\n\x08RESPONSE\x10\r\x12\x16\n\x12GET_STATE_BY_RANGE\x10\x0e\x12\x14\n\x10GET_QUERY_RESULT\x10\x0f\x12\x14\n\x10QUERY_STATE_NEXT\x10\x10\x12\x15\n\x11QUERY_STATE_CLOSE\x10\x11\x12\r\n\tKEEPALIVE\x10\x12\x12\x17\n\x13GET_HISTORY_FOR_KEY\x10\x13\"*\n\x0cPutStateInfo\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c\"3\n\x0fGetStateByRange\x12\x10\n\x08startKey\x18\x01 \x01(\t\x12\x0e\n\x06\x65ndKey\x18\x02 \x01(\t\"\x1f\n\x0eGetQueryResult\x12\r\n\x05query\x18\x01 \x01(\t\"\x1f\n\x10GetHistoryForKey\x12\x0b\n\x03key\x18\x01 \x01(\t\"\x1c\n\x0eQueryStateNext\x12\n\n\x02id\x18\x01 \x01(\t\"\x1d\n\x0fQueryStateClose\x12\n\n\x02id\x18\x01 \x01(\t\"0\n\x12QueryStateKeyValue\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c\"g\n\x12QueryStateResponse\x12\x33\n\x0fkeys_and_values\x18\x01 \x03(\x0b\x32\x1a.protos.QueryStateKeyValue\x12\x10\n\x08has_more\x18\x02 \x01(\x08\x12\n\n\x02id\x18\x03 \x01(\t2X\n\x10\x43haincodeSupport\x12\x44\n\x08Register\x12\x18.protos.ChaincodeMessage\x1a\x18.protos.ChaincodeMessage\"\x00(\x01\x30\x01\x42O\n\"org.hyperledger.fabric.protos.peerZ)github.com/hyperledger/fabric/protos/peerb\x06proto3')
  ,
  dependencies=[peer_dot_chaincode__event__pb2.DESCRIPTOR,peer_dot_proposal__pb2.DESCRIPTOR,google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)



_CHAINCODEMESSAGE_TYPE = _descriptor.EnumDescriptor(
  name='Type',
  full_name='protos.ChaincodeMessage.Type',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNDEFINED', index=0, number=0,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='REGISTER', index=1, number=1,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='REGISTERED', index=2, number=2,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='INIT', index=3, number=3,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='READY', index=4, number=4,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='TRANSACTION', index=5, number=5,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='COMPLETED', index=6, number=6,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='ERROR', index=7, number=7,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GET_STATE', index=8, number=8,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='PUT_STATE', index=9, number=9,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='DEL_STATE', index=10, number=10,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='INVOKE_CHAINCODE', index=11, number=11,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='RESPONSE', index=12, number=13,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GET_STATE_BY_RANGE', index=13, number=14,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GET_QUERY_RESULT', index=14, number=15,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='QUERY_STATE_NEXT', index=15, number=16,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='QUERY_STATE_CLOSE', index=16, number=17,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='KEEPALIVE', index=17, number=18,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GET_HISTORY_FOR_KEY', index=18, number=19,
      options=None,
      type=None),
  ],
  containing_type=None,
  options=None,
  serialized_start=349,
  serialized_end=676,
)
_sym_db.RegisterEnumDescriptor(_CHAINCODEMESSAGE_TYPE)


_CHAINCODEMESSAGE = _descriptor.Descriptor(
  name='ChaincodeMessage',
  full_name='protos.ChaincodeMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='type', full_name='protos.ChaincodeMessage.type', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='protos.ChaincodeMessage.timestamp', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='payload', full_name='protos.ChaincodeMessage.payload', index=2,
      number=3, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='txid', full_name='protos.ChaincodeMessage.txid', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='proposal', full_name='protos.ChaincodeMessage.proposal', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='chaincode_event', full_name='protos.ChaincodeMessage.chaincode_event', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _CHAINCODEMESSAGE_TYPE,
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=120,
  serialized_end=676,
)


_PUTSTATEINFO = _descriptor.Descriptor(
  name='PutStateInfo',
  full_name='protos.PutStateInfo',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='protos.PutStateInfo.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='value', full_name='protos.PutStateInfo.value', index=1,
      number=2, type=12, cpp_type=9, label=1,
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
  serialized_start=678,
  serialized_end=720,
)


_GETSTATEBYRANGE = _descriptor.Descriptor(
  name='GetStateByRange',
  full_name='protos.GetStateByRange',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='startKey', full_name='protos.GetStateByRange.startKey', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='endKey', full_name='protos.GetStateByRange.endKey', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=722,
  serialized_end=773,
)


_GETQUERYRESULT = _descriptor.Descriptor(
  name='GetQueryResult',
  full_name='protos.GetQueryResult',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='query', full_name='protos.GetQueryResult.query', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=775,
  serialized_end=806,
)


_GETHISTORYFORKEY = _descriptor.Descriptor(
  name='GetHistoryForKey',
  full_name='protos.GetHistoryForKey',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='protos.GetHistoryForKey.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=808,
  serialized_end=839,
)


_QUERYSTATENEXT = _descriptor.Descriptor(
  name='QueryStateNext',
  full_name='protos.QueryStateNext',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='protos.QueryStateNext.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=841,
  serialized_end=869,
)


_QUERYSTATECLOSE = _descriptor.Descriptor(
  name='QueryStateClose',
  full_name='protos.QueryStateClose',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='protos.QueryStateClose.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=871,
  serialized_end=900,
)


_QUERYSTATEKEYVALUE = _descriptor.Descriptor(
  name='QueryStateKeyValue',
  full_name='protos.QueryStateKeyValue',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='protos.QueryStateKeyValue.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='value', full_name='protos.QueryStateKeyValue.value', index=1,
      number=2, type=12, cpp_type=9, label=1,
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
  serialized_start=902,
  serialized_end=950,
)


_QUERYSTATERESPONSE = _descriptor.Descriptor(
  name='QueryStateResponse',
  full_name='protos.QueryStateResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='keys_and_values', full_name='protos.QueryStateResponse.keys_and_values', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='has_more', full_name='protos.QueryStateResponse.has_more', index=1,
      number=2, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='id', full_name='protos.QueryStateResponse.id', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=952,
  serialized_end=1055,
)

_CHAINCODEMESSAGE.fields_by_name['type'].enum_type = _CHAINCODEMESSAGE_TYPE
_CHAINCODEMESSAGE.fields_by_name['timestamp'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_CHAINCODEMESSAGE.fields_by_name['proposal'].message_type = peer_dot_proposal__pb2._PROPOSAL
_CHAINCODEMESSAGE.fields_by_name['chaincode_event'].message_type = peer_dot_chaincode__event__pb2._CHAINCODEEVENT
_CHAINCODEMESSAGE_TYPE.containing_type = _CHAINCODEMESSAGE
_QUERYSTATERESPONSE.fields_by_name['keys_and_values'].message_type = _QUERYSTATEKEYVALUE
DESCRIPTOR.message_types_by_name['ChaincodeMessage'] = _CHAINCODEMESSAGE
DESCRIPTOR.message_types_by_name['PutStateInfo'] = _PUTSTATEINFO
DESCRIPTOR.message_types_by_name['GetStateByRange'] = _GETSTATEBYRANGE
DESCRIPTOR.message_types_by_name['GetQueryResult'] = _GETQUERYRESULT
DESCRIPTOR.message_types_by_name['GetHistoryForKey'] = _GETHISTORYFORKEY
DESCRIPTOR.message_types_by_name['QueryStateNext'] = _QUERYSTATENEXT
DESCRIPTOR.message_types_by_name['QueryStateClose'] = _QUERYSTATECLOSE
DESCRIPTOR.message_types_by_name['QueryStateKeyValue'] = _QUERYSTATEKEYVALUE
DESCRIPTOR.message_types_by_name['QueryStateResponse'] = _QUERYSTATERESPONSE

ChaincodeMessage = _reflection.GeneratedProtocolMessageType('ChaincodeMessage', (_message.Message,), dict(
  DESCRIPTOR = _CHAINCODEMESSAGE,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.ChaincodeMessage)
  ))
_sym_db.RegisterMessage(ChaincodeMessage)

PutStateInfo = _reflection.GeneratedProtocolMessageType('PutStateInfo', (_message.Message,), dict(
  DESCRIPTOR = _PUTSTATEINFO,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.PutStateInfo)
  ))
_sym_db.RegisterMessage(PutStateInfo)

GetStateByRange = _reflection.GeneratedProtocolMessageType('GetStateByRange', (_message.Message,), dict(
  DESCRIPTOR = _GETSTATEBYRANGE,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.GetStateByRange)
  ))
_sym_db.RegisterMessage(GetStateByRange)

GetQueryResult = _reflection.GeneratedProtocolMessageType('GetQueryResult', (_message.Message,), dict(
  DESCRIPTOR = _GETQUERYRESULT,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.GetQueryResult)
  ))
_sym_db.RegisterMessage(GetQueryResult)

GetHistoryForKey = _reflection.GeneratedProtocolMessageType('GetHistoryForKey', (_message.Message,), dict(
  DESCRIPTOR = _GETHISTORYFORKEY,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.GetHistoryForKey)
  ))
_sym_db.RegisterMessage(GetHistoryForKey)

QueryStateNext = _reflection.GeneratedProtocolMessageType('QueryStateNext', (_message.Message,), dict(
  DESCRIPTOR = _QUERYSTATENEXT,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.QueryStateNext)
  ))
_sym_db.RegisterMessage(QueryStateNext)

QueryStateClose = _reflection.GeneratedProtocolMessageType('QueryStateClose', (_message.Message,), dict(
  DESCRIPTOR = _QUERYSTATECLOSE,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.QueryStateClose)
  ))
_sym_db.RegisterMessage(QueryStateClose)

QueryStateKeyValue = _reflection.GeneratedProtocolMessageType('QueryStateKeyValue', (_message.Message,), dict(
  DESCRIPTOR = _QUERYSTATEKEYVALUE,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.QueryStateKeyValue)
  ))
_sym_db.RegisterMessage(QueryStateKeyValue)

QueryStateResponse = _reflection.GeneratedProtocolMessageType('QueryStateResponse', (_message.Message,), dict(
  DESCRIPTOR = _QUERYSTATERESPONSE,
  __module__ = 'peer.chaincode_shim_pb2'
  # @@protoc_insertion_point(class_scope:protos.QueryStateResponse)
  ))
_sym_db.RegisterMessage(QueryStateResponse)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('\n\"org.hyperledger.fabric.protos.peerZ)github.com/hyperledger/fabric/protos/peer'))
try:
  # THESE ELEMENTS WILL BE DEPRECATED.
  # Please use the generated *_pb2_grpc.py files instead.
  import grpc
  from grpc.framework.common import cardinality
  from grpc.framework.interfaces.face import utilities as face_utilities
  from grpc.beta import implementations as beta_implementations
  from grpc.beta import interfaces as beta_interfaces


  class ChaincodeSupportStub(object):
    """Interface that provides support to chaincode execution. ChaincodeContext
    provides the context necessary for the server to respond appropriately.
    """

    def __init__(self, channel):
      """Constructor.

      Args:
        channel: A grpc.Channel.
      """
      self.Register = channel.stream_stream(
          '/protos.ChaincodeSupport/Register',
          request_serializer=ChaincodeMessage.SerializeToString,
          response_deserializer=ChaincodeMessage.FromString,
          )


  class ChaincodeSupportServicer(object):
    """Interface that provides support to chaincode execution. ChaincodeContext
    provides the context necessary for the server to respond appropriately.
    """

    def Register(self, request_iterator, context):
      context.set_code(grpc.StatusCode.UNIMPLEMENTED)
      context.set_details('Method not implemented!')
      raise NotImplementedError('Method not implemented!')


  def add_ChaincodeSupportServicer_to_server(servicer, server):
    rpc_method_handlers = {
        'Register': grpc.stream_stream_rpc_method_handler(
            servicer.Register,
            request_deserializer=ChaincodeMessage.FromString,
            response_serializer=ChaincodeMessage.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        'protos.ChaincodeSupport', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


  class BetaChaincodeSupportServicer(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """Interface that provides support to chaincode execution. ChaincodeContext
    provides the context necessary for the server to respond appropriately.
    """
    def Register(self, request_iterator, context):
      context.code(beta_interfaces.StatusCode.UNIMPLEMENTED)


  class BetaChaincodeSupportStub(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """Interface that provides support to chaincode execution. ChaincodeContext
    provides the context necessary for the server to respond appropriately.
    """
    def Register(self, request_iterator, timeout, metadata=None, with_call=False, protocol_options=None):
      raise NotImplementedError()


  def beta_create_ChaincodeSupport_server(servicer, pool=None, pool_size=None, default_timeout=None, maximum_timeout=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_deserializers = {
      ('protos.ChaincodeSupport', 'Register'): ChaincodeMessage.FromString,
    }
    response_serializers = {
      ('protos.ChaincodeSupport', 'Register'): ChaincodeMessage.SerializeToString,
    }
    method_implementations = {
      ('protos.ChaincodeSupport', 'Register'): face_utilities.stream_stream_inline(servicer.Register),
    }
    server_options = beta_implementations.server_options(request_deserializers=request_deserializers, response_serializers=response_serializers, thread_pool=pool, thread_pool_size=pool_size, default_timeout=default_timeout, maximum_timeout=maximum_timeout)
    return beta_implementations.server(method_implementations, options=server_options)


  def beta_create_ChaincodeSupport_stub(channel, host=None, metadata_transformer=None, pool=None, pool_size=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_serializers = {
      ('protos.ChaincodeSupport', 'Register'): ChaincodeMessage.SerializeToString,
    }
    response_deserializers = {
      ('protos.ChaincodeSupport', 'Register'): ChaincodeMessage.FromString,
    }
    cardinalities = {
      'Register': cardinality.Cardinality.STREAM_STREAM,
    }
    stub_options = beta_implementations.stub_options(host=host, metadata_transformer=metadata_transformer, request_serializers=request_serializers, response_deserializers=response_deserializers, thread_pool=pool, thread_pool_size=pool_size)
    return beta_implementations.dynamic_stub(channel, 'protos.ChaincodeSupport', cardinalities, options=stub_options)
except ImportError:
  pass
# @@protoc_insertion_point(module_scope)
