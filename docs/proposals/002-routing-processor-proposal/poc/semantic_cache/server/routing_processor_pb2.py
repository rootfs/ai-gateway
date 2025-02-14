# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: routing-processor.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'routing-processor.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x17routing-processor.proto\x12\x1arouting_processor.v1alpha1\"\xd0\x01\n\x05Value\x12\x16\n\x0cstring_value\x18\x01 \x01(\tH\x00\x12\x16\n\x0cnumber_value\x18\x02 \x01(\x01H\x00\x12\x14\n\nbool_value\x18\x03 \x01(\x08H\x00\x12<\n\x0cstruct_value\x18\x04 \x01(\x0b\x32$.routing_processor.v1alpha1.MetadataH\x00\x12;\n\nlist_value\x18\x05 \x01(\x0b\x32%.routing_processor.v1alpha1.ListValueH\x00\x42\x06\n\x04kind\">\n\tListValue\x12\x31\n\x06values\x18\x01 \x03(\x0b\x32!.routing_processor.v1alpha1.Value\"\x9e\x01\n\x08Metadata\x12@\n\x06\x66ields\x18\x01 \x03(\x0b\x32\x30.routing_processor.v1alpha1.Metadata.FieldsEntry\x1aP\n\x0b\x46ieldsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x30\n\x05value\x18\x02 \x01(\x0b\x32!.routing_processor.v1alpha1.Value:\x02\x38\x01\"$\n\x06Header\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t\"O\n\x05Usage\x12\x15\n\rprompt_tokens\x18\x01 \x01(\x05\x12\x19\n\x11\x63ompletion_tokens\x18\x02 \x01(\x05\x12\x14\n\x0ctotal_tokens\x18\x03 \x01(\x05\"\x9d\x01\n\x07Message\x12\x0c\n\x04role\x18\x01 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x02 \x01(\t\x12\x38\n\ntool_calls\x18\x03 \x03(\x0b\x32$.routing_processor.v1alpha1.ToolCall\x12\x14\n\x0ctool_call_id\x18\x04 \x01(\t\x12\x0c\n\x04name\x18\x05 \x01(\t\x12\x15\n\rfunction_call\x18\x06 \x01(\t\"I\n\x08ToolCall\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x10\n\x08\x66unction\x18\x03 \x01(\t\x12\x11\n\targuments\x18\x04 \x01(\t\"\xa7\x01\n\x0f\x43\x61\x63hingBehavior\x12!\n\x14similarity_threshold\x18\x01 \x01(\x02H\x00\x88\x01\x01\x12\x1b\n\x0emax_cache_size\x18\x02 \x01(\x05H\x01\x88\x01\x01\x12\x18\n\x0bttl_seconds\x18\x03 \x01(\x05H\x02\x88\x01\x01\x42\x17\n\x15_similarity_thresholdB\x11\n\x0f_max_cache_sizeB\x0e\n\x0c_ttl_seconds\"\x8e\x02\n\x0eProcessingMode\x12\x1e\n\x16semantic_cache_enabled\x18\x01 \x01(\x08\x12\x1f\n\x17model_selection_enabled\x18\x02 \x01(\x08\x12\x14\n\x07timeout\x18\x03 \x01(\rH\x00\x88\x01\x01\x12J\n\x10\x63\x61\x63hing_behavior\x18\x04 \x01(\x0b\x32+.routing_processor.v1alpha1.CachingBehaviorH\x01\x88\x01\x01\x12 \n\x13\x61llow_mode_override\x18\x05 \x01(\x08H\x02\x88\x01\x01\x42\n\n\x08_timeoutB\x13\n\x11_caching_behaviorB\x16\n\x14_allow_mode_override\"*\n\tAPISchema\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0f\n\x07version\x18\x02 \x01(\t\"h\n\x0e\x42\x61\x63kendMetrics\x12\x18\n\x0b\x61vg_latency\x18\x01 \x01(\x02H\x00\x88\x01\x01\x12\x1a\n\rcost_per_unit\x18\x02 \x01(\x02H\x01\x88\x01\x01\x42\x0e\n\x0c_avg_latencyB\x10\n\x0e_cost_per_unit\"\xe6\x01\n\x07\x42\x61\x63kend\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0e\n\x06weight\x18\x02 \x01(\x05\x12\x35\n\x06schema\x18\x03 \x01(\x0b\x32%.routing_processor.v1alpha1.APISchema\x12\x38\n\nproperties\x18\x04 \x01(\x0b\x32$.routing_processor.v1alpha1.Metadata\x12@\n\x07metrics\x18\x05 \x01(\x0b\x32*.routing_processor.v1alpha1.BackendMetricsH\x00\x88\x01\x01\x42\n\n\x08_metrics\"`\n\tCacheInfo\x12\x11\n\tcache_hit\x18\x01 \x01(\x08\x12\x18\n\x10similarity_score\x18\x02 \x01(\x02\x12\x19\n\x11\x63\x61\x63hed_request_id\x18\x03 \x01(\t\x12\x0b\n\x03\x61ge\x18\x04 \x01(\r\"W\n\x11HeaderValueOption\x12\x32\n\x06header\x18\x01 \x01(\x0b\x32\".routing_processor.v1alpha1.Header\x12\x0e\n\x06\x61ppend\x18\x02 \x01(\x08\"l\n\x0eHeaderMutation\x12\x42\n\x0bset_headers\x18\x01 \x03(\x0b\x32-.routing_processor.v1alpha1.HeaderValueOption\x12\x16\n\x0eremove_headers\x18\x02 \x03(\t\"\x1c\n\x0c\x42odyMutation\x12\x0c\n\x04\x62ody\x18\x01 \x01(\x0c\"\xa7\x01\n\x0fRequestMetadata\x12\x12\n\nrequest_id\x18\x01 \x01(\t\x12\x12\n\nmodel_name\x18\x02 \x01(\t\x12\x18\n\x10response_latency\x18\x03 \x01(\x02\x12\x13\n\x0btoken_count\x18\x04 \x01(\x04\x12\x13\n\x0b\x63\x61\x63he_score\x18\x05 \x01(\x02\x12\x18\n\x0bstatus_code\x18\x06 \x01(\x05H\x00\x88\x01\x01\x42\x0e\n\x0c_status_code\"\x9a\x01\n\x0eRequestHeaders\x12\x33\n\x07headers\x18\x01 \x03(\x0b\x32\".routing_processor.v1alpha1.Header\x12?\n\x12\x61vailable_backends\x18\x02 \x03(\x0b\x32#.routing_processor.v1alpha1.Backend\x12\x12\n\nrequest_id\x18\x03 \x01(\t\"\xb0\x01\n\x0bRequestBody\x12\x0c\n\x04\x62ody\x18\x01 \x01(\x0c\x12\x15\n\rend_of_stream\x18\x02 \x01(\x08\x12=\n\x0cmessage_type\x18\x03 \x01(\x0e\x32\'.routing_processor.v1alpha1.MessageType\x12=\n\x08metadata\x18\x04 \x01(\x0b\x32+.routing_processor.v1alpha1.RequestMetadata\"\xdb\x01\n\x11ProcessingRequest\x12=\n\x07headers\x18\x01 \x01(\x0b\x32*.routing_processor.v1alpha1.RequestHeadersH\x00\x12\x37\n\x04\x62ody\x18\x02 \x01(\x0b\x32\'.routing_processor.v1alpha1.RequestBodyH\x00\x12\x43\n\x0fprocessing_mode\x18\x03 \x01(\x0b\x32*.routing_processor.v1alpha1.ProcessingModeB\t\n\x07request\"\x99\x02\n\x0e\x43ommonResponse\x12\x43\n\x0fheader_mutation\x18\x01 \x01(\x0b\x32*.routing_processor.v1alpha1.HeaderMutation\x12?\n\rbody_mutation\x18\x02 \x01(\x0b\x32(.routing_processor.v1alpha1.BodyMutation\x12>\n\ncache_info\x18\x03 \x01(\x0b\x32%.routing_processor.v1alpha1.CacheInfoH\x00\x88\x01\x01\x12\x1d\n\x10selected_backend\x18\x04 \x01(\tH\x01\x88\x01\x01\x42\r\n\x0b_cache_infoB\x13\n\x11_selected_backend\"k\n\x11ImmediateResponse\x12\x13\n\x0bstatus_code\x18\x01 \x01(\x05\x12\x33\n\x07headers\x18\x02 \x03(\x0b\x32\".routing_processor.v1alpha1.Header\x12\x0c\n\x04\x62ody\x18\x03 \x01(\x0c\"Y\n\x19MessageProcessingResponse\x12<\n\x08response\x18\x01 \x01(\x0b\x32*.routing_processor.v1alpha1.CommonResponse\"\x82\x02\n\x12ProcessingResponse\x12K\n\x12immediate_response\x18\x01 \x01(\x0b\x32-.routing_processor.v1alpha1.ImmediateResponseH\x00\x12S\n\x12message_processing\x18\x02 \x01(\x0b\x32\x35.routing_processor.v1alpha1.MessageProcessingResponseH\x00\x12>\n\x10\x64ynamic_metadata\x18\x03 \x01(\x0b\x32$.routing_processor.v1alpha1.MetadataB\n\n\x08response\"s\n\rSearchRequest\x12\x35\n\x08messages\x18\x01 \x03(\x0b\x32#.routing_processor.v1alpha1.Message\x12\x1c\n\x14similarity_threshold\x18\x02 \x01(\x02\x12\r\n\x05model\x18\x03 \x01(\t\"\xab\x01\n\x0eSearchResponse\x12\r\n\x05\x66ound\x18\x01 \x01(\x08\x12>\n\x11response_messages\x18\x02 \x03(\x0b\x32#.routing_processor.v1alpha1.Message\x12\x18\n\x10similarity_score\x18\x03 \x01(\x02\x12\x30\n\x05usage\x18\x04 \x01(\x0b\x32!.routing_processor.v1alpha1.Usage\"\xdf\x01\n\x10StoreChatRequest\x12=\n\x10request_messages\x18\x01 \x03(\x0b\x32#.routing_processor.v1alpha1.Message\x12>\n\x11response_messages\x18\x02 \x03(\x0b\x32#.routing_processor.v1alpha1.Message\x12\r\n\x05model\x18\x03 \x01(\t\x12\x0b\n\x03ttl\x18\x04 \x01(\t\x12\x30\n\x05usage\x18\x05 \x01(\x0b\x32!.routing_processor.v1alpha1.Usage\"3\n\x11StoreChatResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\r\n\x05\x65rror\x18\x02 \x01(\t\"\x15\n\x13\x43\x61pabilitiesRequest\"\x81\x01\n\x14\x43\x61pabilitiesResponse\x12 \n\x18semantic_cache_supported\x18\x01 \x01(\x08\x12!\n\x19model_selection_supported\x18\x02 \x01(\x08\x12$\n\x1cimmediate_response_supported\x18\x03 \x01(\x08*Z\n\x0bMessageType\x12\x0b\n\x07UNKNOWN\x10\x00\x12\x12\n\x0e\x43LIENT_REQUEST\x10\x01\x12\x15\n\x11UPSTREAM_RESPONSE\x10\x02\x12\x13\n\x0f\x43\x41\x43HED_RESPONSE\x10\x03\x32\xfe\x01\n\x10RoutingProcessor\x12t\n\x0f\x45xternalProcess\x12-.routing_processor.v1alpha1.ProcessingRequest\x1a..routing_processor.v1alpha1.ProcessingResponse(\x01\x30\x01\x12t\n\x0fGetCapabilities\x12/.routing_processor.v1alpha1.CapabilitiesRequest\x1a\x30.routing_processor.v1alpha1.CapabilitiesResponse2\xea\x01\n\x14SemanticCacheService\x12\x66\n\x0bSearchCache\x12).routing_processor.v1alpha1.SearchRequest\x1a*.routing_processor.v1alpha1.SearchResponse\"\x00\x12j\n\tStoreChat\x12,.routing_processor.v1alpha1.StoreChatRequest\x1a-.routing_processor.v1alpha1.StoreChatResponse\"\x00\x42\x16Z\x14routing_processor/gob\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'routing_processor_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\024routing_processor/go'
  _globals['_METADATA_FIELDSENTRY']._loaded_options = None
  _globals['_METADATA_FIELDSENTRY']._serialized_options = b'8\001'
  _globals['_MESSAGETYPE']._serialized_start=4196
  _globals['_MESSAGETYPE']._serialized_end=4286
  _globals['_VALUE']._serialized_start=56
  _globals['_VALUE']._serialized_end=264
  _globals['_LISTVALUE']._serialized_start=266
  _globals['_LISTVALUE']._serialized_end=328
  _globals['_METADATA']._serialized_start=331
  _globals['_METADATA']._serialized_end=489
  _globals['_METADATA_FIELDSENTRY']._serialized_start=409
  _globals['_METADATA_FIELDSENTRY']._serialized_end=489
  _globals['_HEADER']._serialized_start=491
  _globals['_HEADER']._serialized_end=527
  _globals['_USAGE']._serialized_start=529
  _globals['_USAGE']._serialized_end=608
  _globals['_MESSAGE']._serialized_start=611
  _globals['_MESSAGE']._serialized_end=768
  _globals['_TOOLCALL']._serialized_start=770
  _globals['_TOOLCALL']._serialized_end=843
  _globals['_CACHINGBEHAVIOR']._serialized_start=846
  _globals['_CACHINGBEHAVIOR']._serialized_end=1013
  _globals['_PROCESSINGMODE']._serialized_start=1016
  _globals['_PROCESSINGMODE']._serialized_end=1286
  _globals['_APISCHEMA']._serialized_start=1288
  _globals['_APISCHEMA']._serialized_end=1330
  _globals['_BACKENDMETRICS']._serialized_start=1332
  _globals['_BACKENDMETRICS']._serialized_end=1436
  _globals['_BACKEND']._serialized_start=1439
  _globals['_BACKEND']._serialized_end=1669
  _globals['_CACHEINFO']._serialized_start=1671
  _globals['_CACHEINFO']._serialized_end=1767
  _globals['_HEADERVALUEOPTION']._serialized_start=1769
  _globals['_HEADERVALUEOPTION']._serialized_end=1856
  _globals['_HEADERMUTATION']._serialized_start=1858
  _globals['_HEADERMUTATION']._serialized_end=1966
  _globals['_BODYMUTATION']._serialized_start=1968
  _globals['_BODYMUTATION']._serialized_end=1996
  _globals['_REQUESTMETADATA']._serialized_start=1999
  _globals['_REQUESTMETADATA']._serialized_end=2166
  _globals['_REQUESTHEADERS']._serialized_start=2169
  _globals['_REQUESTHEADERS']._serialized_end=2323
  _globals['_REQUESTBODY']._serialized_start=2326
  _globals['_REQUESTBODY']._serialized_end=2502
  _globals['_PROCESSINGREQUEST']._serialized_start=2505
  _globals['_PROCESSINGREQUEST']._serialized_end=2724
  _globals['_COMMONRESPONSE']._serialized_start=2727
  _globals['_COMMONRESPONSE']._serialized_end=3008
  _globals['_IMMEDIATERESPONSE']._serialized_start=3010
  _globals['_IMMEDIATERESPONSE']._serialized_end=3117
  _globals['_MESSAGEPROCESSINGRESPONSE']._serialized_start=3119
  _globals['_MESSAGEPROCESSINGRESPONSE']._serialized_end=3208
  _globals['_PROCESSINGRESPONSE']._serialized_start=3211
  _globals['_PROCESSINGRESPONSE']._serialized_end=3469
  _globals['_SEARCHREQUEST']._serialized_start=3471
  _globals['_SEARCHREQUEST']._serialized_end=3586
  _globals['_SEARCHRESPONSE']._serialized_start=3589
  _globals['_SEARCHRESPONSE']._serialized_end=3760
  _globals['_STORECHATREQUEST']._serialized_start=3763
  _globals['_STORECHATREQUEST']._serialized_end=3986
  _globals['_STORECHATRESPONSE']._serialized_start=3988
  _globals['_STORECHATRESPONSE']._serialized_end=4039
  _globals['_CAPABILITIESREQUEST']._serialized_start=4041
  _globals['_CAPABILITIESREQUEST']._serialized_end=4062
  _globals['_CAPABILITIESRESPONSE']._serialized_start=4065
  _globals['_CAPABILITIESRESPONSE']._serialized_end=4194
  _globals['_ROUTINGPROCESSOR']._serialized_start=4289
  _globals['_ROUTINGPROCESSOR']._serialized_end=4543
  _globals['_SEMANTICCACHESERVICE']._serialized_start=4546
  _globals['_SEMANTICCACHESERVICE']._serialized_end=4780
# @@protoc_insertion_point(module_scope)
