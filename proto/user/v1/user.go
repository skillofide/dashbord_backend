package userv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Message Types ────────────────────────────────────────────────────────────

type VerifyUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyUserResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type CreateOrUpdateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateOrUpdateUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ─── Profile Message Types ────────────────────────────────────────────────────

// UserProfile holds all profile fields mirroring the frontend ProfilePage.
type UserProfile struct {
	UserID string `json:"user_id"`

	// Personal Info
	Gender     string `json:"gender"`
	Dob        string `json:"dob"`
	Whatsapp   string `json:"whatsapp"`
	Phone      string `json:"phone"`
	Experience string `json:"experience"`

	// Generic Details
	WorkExperience         string   `json:"work_experience"`
	CareerGap              string   `json:"career_gap"`
	CurrentState           string   `json:"current_state"`
	CurrentCity            string   `json:"current_city"`
	PreferredLocations     []string `json:"preferred_locations"`
	GithubLink             string   `json:"github_link"`
	LinkedinLink           string   `json:"linkedin_link"`
	IsWorkingProfessional  bool     `json:"is_working_professional"`
	ResumeName             string   `json:"resume_name"`

	// 10th Grade
	Edu10SchoolName    string `json:"edu10_school_name"`
	Edu10YearOfPassout string `json:"edu10_year_of_passout"`
	Edu10MarksPercent  string `json:"edu10_marks_percent"`

	// 12th / PUC / Intermediate / Diploma
	Edu12SchoolName    string `json:"edu12_school_name"`
	Edu12YearOfPassout string `json:"edu12_year_of_passout"`
	Edu12MarksPercent  string `json:"edu12_marks_percent"`

	// UG Detail
	UGUniversityRollNo string `json:"ug_university_roll_no"`
	UGCollegeName      string `json:"ug_college_name"`
	UGCourseName       string `json:"ug_course_name"`
	UGBranch           string `json:"ug_branch"`
	UGYearOfPassout    string `json:"ug_year_of_passout"`
	UGMarksPercent     string `json:"ug_marks_percent"`
	UGCGPA             string `json:"ug_cgpa"`
	UGActiveBacklogs   string `json:"ug_active_backlogs"`

	// PG Detail
	PGHasCertificate bool `json:"pg_has_certificate"`
}

type GetProfileRequest struct {
	UserID string `json:"user_id"`
}

type GetProfileResponse struct {
	Profile *UserProfile `json:"profile"`
}

type UpsertProfileRequest struct {
	Profile *UserProfile `json:"profile"`
}

type UpsertProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ─── Server Interface ─────────────────────────────────────────────────────────

type UserServiceServer interface {
	VerifyUser(context.Context, *VerifyUserRequest) (*VerifyUserResponse, error)
	CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error)
	GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
	UpsertProfile(context.Context, *UpsertProfileRequest) (*UpsertProfileResponse, error)
}

type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) VerifyUser(context.Context, *VerifyUserRequest) (*VerifyUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyUser not implemented")
}

func (UnimplementedUserServiceServer) CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateUser not implemented")
}

func (UnimplementedUserServiceServer) GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}

func (UnimplementedUserServiceServer) UpsertProfile(context.Context, *UpsertProfileRequest) (*UpsertProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertProfile not implemented")
}

// ─── Client Interface & Implementation ───────────────────────────────────────

type UserServiceClient interface {
	VerifyUser(ctx context.Context, in *VerifyUserRequest, opts ...grpc.CallOption) (*VerifyUserResponse, error)
	CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error)
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error)
	UpsertProfile(ctx context.Context, in *UpsertProfileRequest, opts ...grpc.CallOption) (*UpsertProfileResponse, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) VerifyUser(ctx context.Context, in *VerifyUserRequest, opts ...grpc.CallOption) (*VerifyUserResponse, error) {
	out := new(VerifyUserResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/VerifyUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error) {
	out := new(CreateOrUpdateUserResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/CreateOrUpdateUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error) {
	out := new(GetProfileResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/GetProfile", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpsertProfile(ctx context.Context, in *UpsertProfileRequest, opts ...grpc.CallOption) (*UpsertProfileResponse, error) {
	out := new(UpsertProfileResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/UpsertProfile", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Service Registration & Descriptor ───────────────────────────────────────

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "VerifyUser", Handler: _UserService_VerifyUser_Handler},
		{MethodName: "CreateOrUpdateUser", Handler: _UserService_CreateOrUpdateUser_Handler},
		{MethodName: "GetProfile", Handler: _UserService_GetProfile_Handler},
		{MethodName: "UpsertProfile", Handler: _UserService_UpsertProfile_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _UserService_VerifyUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyUser(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/VerifyUser"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).VerifyUser(ctx, req.(*VerifyUserRequest))
		})
}

func _UserService_CreateOrUpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateOrUpdateUser(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/CreateOrUpdateUser"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).CreateOrUpdateUser(ctx, req.(*CreateOrUpdateUserRequest))
		})
}

func _UserService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetProfile(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/GetProfile"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).GetProfile(ctx, req.(*GetProfileRequest))
		})
}

func _UserService_UpsertProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpsertProfile(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/UpsertProfile"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).UpsertProfile(ctx, req.(*UpsertProfileRequest))
		})
}

