package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IngestRequest struct {
	state protoimpl.MessageState
	sizeCache protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdRequisicao  string
	MapperVersion string
	WebhookUrl    string
	TaxRate       float64
	CsvData       []byte
}

func (x *IngestRequest) Reset() {
	*x = IngestRequest{}
}

func (x *IngestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngestRequest) ProtoMessage() {}

func (x *IngestRequest) GetIdRequisicao() string {
	if x != nil {
		return x.IdRequisicao
	}
	return ""
}

func (x *IngestRequest) GetMapperVersion() string {
	if x != nil {
		return x.MapperVersion
	}
	return ""
}

func (x *IngestRequest) GetWebhookUrl() string {
	if x != nil {
		return x.WebhookUrl
	}
	return ""
}

func (x *IngestRequest) GetTaxRate() float64 {
	if x != nil {
		return x.TaxRate
	}
	return 0
}

func (x *IngestRequest) GetCsvData() []byte {
	if x != nil {
		return x.CsvData
	}
	return nil
}

type IngestResponse struct {
	state protoimpl.MessageState
	sizeCache protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdRequisicao string
	Status       string
	DocumentId   int64
	Error        string
}

func (x *IngestResponse) Reset() {
	*x = IngestResponse{}
}

func (x *IngestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngestResponse) ProtoMessage() {}

func (x *IngestResponse) GetIdRequisicao() string {
	if x != nil {
		return x.IdRequisicao
	}
	return ""
}

func (x *IngestResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *IngestResponse) GetDocumentId() int64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *IngestResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type CategoryQuery struct {
	state protoimpl.MessageState
	Category string
}

func (x *CategoryQuery) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

type SalesResult struct {
	state protoimpl.MessageState
	Category    string
	TotalVendas float64
}

func (x *SalesResult) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *SalesResult) GetTotalVendas() float64 {
	if x != nil {
		return x.TotalVendas
	}
	return 0
}

type RegionQuery struct {
	state protoimpl.MessageState
	Region string
}

func (x *RegionQuery) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

type ProfitResult struct {
	state protoimpl.MessageState
	Region      string
	LucroTotal  float64
}

func (x *ProfitResult) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ProfitResult) GetLucroTotal() float64 {
	if x != nil {
		return x.LucroTotal
	}
	return 0
}

type LimitQuery struct {
	state protoimpl.MessageState
	Limit int32
}

func (x *LimitQuery) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type LossOrder struct {
	state protoimpl.MessageState
	OrderId    string
	LucroTotal float64
}

func (x *LossOrder) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *LossOrder) GetLucroTotal() float64 {
	if x != nil {
		return x.LucroTotal
	}
	return 0
}

type LossOrdersResult struct {
	state protoimpl.MessageState
	Results []*LossOrder
}

func (x *LossOrdersResult) GetResults() []*LossOrder {
	if x != nil {
		return x.Results
	}
	return nil
}

// ===== gRPC Service Definitions =====

const (
	IngestService_ServiceDesc_name_IngestCSV = "IngestCSV"
	QueryService_ServiceDesc_name_GetSalesByCategory = "GetSalesByCategory"
	QueryService_ServiceDesc_name_GetProfitByRegion = "GetProfitByRegion"
	QueryService_ServiceDesc_name_GetLossOrders = "GetLossOrders"
)

type UnimplementedIngestServiceServer struct {}

func (UnimplementedIngestServiceServer) IngestCSV(context.Context, *IngestRequest) (*IngestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IngestCSV not implemented")
}

type UnimplementedQueryServiceServer struct {}

func (UnimplementedQueryServiceServer) GetSalesByCategory(context.Context, *CategoryQuery) (*SalesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSalesByCategory not implemented")
}

func (UnimplementedQueryServiceServer) GetProfitByRegion(context.Context, *RegionQuery) (*ProfitResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfitByRegion not implemented")
}

func (UnimplementedQueryServiceServer) GetLossOrders(context.Context, *LimitQuery) (*LossOrdersResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLossOrders not implemented")
}

type IngestServiceClient interface {
	IngestCSV(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error)
}

type ingestServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIngestServiceClient(cc grpc.ClientConnInterface) IngestServiceClient {
	return &ingestServiceClient{cc}
}

func (c *ingestServiceClient) IngestCSV(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error) {
	out := new(IngestResponse)
	err := c.cc.Invoke(ctx, "/xmlservice.IngestService/IngestCSV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type QueryServiceClient interface {
	GetSalesByCategory(ctx context.Context, in *CategoryQuery, opts ...grpc.CallOption) (*SalesResult, error)
	GetProfitByRegion(ctx context.Context, in *RegionQuery, opts ...grpc.CallOption) (*ProfitResult, error)
	GetLossOrders(ctx context.Context, in *LimitQuery, opts ...grpc.CallOption) (*LossOrdersResult, error)
}

type queryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryServiceClient(cc grpc.ClientConnInterface) QueryServiceClient {
	return &queryServiceClient{cc}
}

func (c *queryServiceClient) GetSalesByCategory(ctx context.Context, in *CategoryQuery, opts ...grpc.CallOption) (*SalesResult, error) {
	out := new(SalesResult)
	err := c.cc.Invoke(ctx, "/xmlservice.QueryService/GetSalesByCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryServiceClient) GetProfitByRegion(ctx context.Context, in *RegionQuery, opts ...grpc.CallOption) (*ProfitResult, error) {
	out := new(ProfitResult)
	err := c.cc.Invoke(ctx, "/xmlservice.QueryService/GetProfitByRegion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryServiceClient) GetLossOrders(ctx context.Context, in *LimitQuery, opts ...grpc.CallOption) (*LossOrdersResult, error) {
	out := new(LossOrdersResult)
	err := c.cc.Invoke(ctx, "/xmlservice.QueryService/GetLossOrders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func RegisterIngestServiceServer(s grpc.ServiceRegistrar, srv IngestServiceServer) {
	s.RegisterService(&IngestService_ServiceDesc, srv)
}

func RegisterQueryServiceServer(s grpc.ServiceRegistrar, srv QueryServiceServer) {
	s.RegisterService(&QueryService_ServiceDesc, srv)
}

var IngestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xmlservice.IngestService",
	HandlerType: (*IngestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IngestCSV",
			Handler:    _IngestService_IngestCSV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ingest.proto",
}

var QueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xmlservice.QueryService",
	HandlerType: (*QueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSalesByCategory",
			Handler:    _QueryService_GetSalesByCategory_Handler,
		},
		{
			MethodName: "GetProfitByRegion",
			Handler:    _QueryService_GetProfitByRegion_Handler,
		},
		{
			MethodName: "GetLossOrders",
			Handler:    _QueryService_GetLossOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "query.proto",
}

type IngestServiceServer interface {
	IngestCSV(context.Context, *IngestRequest) (*IngestResponse, error)
}

type QueryServiceServer interface {
	GetSalesByCategory(context.Context, *CategoryQuery) (*SalesResult, error)
	GetProfitByRegion(context.Context, *RegionQuery) (*ProfitResult, error)
	GetLossOrders(context.Context, *LimitQuery) (*LossOrdersResult, error)
}

func _IngestService_IngestCSV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestServiceServer).IngestCSV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xmlservice.IngestService/IngestCSV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestServiceServer).IngestCSV(ctx, req.(*IngestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueryService_GetSalesByCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CategoryQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServiceServer).GetSalesByCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xmlservice.QueryService/GetSalesByCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServiceServer).GetSalesByCategory(ctx, req.(*CategoryQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueryService_GetProfitByRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegionQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServiceServer).GetProfitByRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xmlservice.QueryService/GetProfitByRegion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServiceServer).GetProfitByRegion(ctx, req.(*RegionQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueryService_GetLossOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LimitQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServiceServer).GetLossOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xmlservice.QueryService/GetLossOrders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServiceServer).GetLossOrders(ctx, req.(*LimitQuery))
	}
	return interceptor(ctx, in, info, handler)
}
