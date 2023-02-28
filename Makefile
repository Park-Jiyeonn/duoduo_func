kitex_gen_user:
	kitex -module simple_tiktok -service simple_tiktok idl/base.thrift # execute in the project root directory
	kitex -module simple_tiktok -service simple_tiktok idl/interact.thrift