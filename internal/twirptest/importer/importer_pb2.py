# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: importer.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from github.com.twitchtv.twirp.internal.twirptest.importable import importable_pb2 as github_dot_com_dot_twitchtv_dot_twirp_dot_internal_dot_twirptest_dot_importable_dot_importable__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='importer.proto',
  package='twirp.internal.twirptest.importer',
  syntax='proto3',
  serialized_pb=_b('\n\x0eimporter.proto\x12!twirp.internal.twirptest.importer\x1aHgithub.com/twitchtv/twirp/internal/twirptest/importable/importable.proto2b\n\x04Svc2\x12Z\n\x04Send\x12(.twirp.internal.twirptest.importable.Msg\x1a(.twirp.internal.twirptest.importable.MsgB\nZ\x08importerb\x06proto3')
  ,
  dependencies=[github_dot_com_dot_twitchtv_dot_twirp_dot_internal_dot_twirptest_dot_importable_dot_importable__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('Z\010importer'))
# @@protoc_insertion_point(module_scope)