package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	vminterfaces "github.com/flashhhhh/Virtual-Machines-Management/custom-controller/vm_interfaces"
)

type AWSProvider struct {
	client *ec2.Client
}

func NewAWSProvider() (*AWSProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, errors.New("unable to load AWS SDK config: " + err.Error())
	}

	return &AWSProvider{
		client: ec2.NewFromConfig(cfg),
	}, nil
}

func (a *AWSProvider) CreateVM(cfg vminterfaces.VirtualMachineConfig) (*vminterfaces.VirtualMachine, error) {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(cfg.Image),
		InstanceType: types.InstanceType(cfg.Size),
		KeyName:      aws.String(cfg.SSHKeyIDs[0]),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     aws.String(cfg.SubnetID),
		SecurityGroupIds: cfg.SecurityGroupIDs,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{Key: aws.String("Name"), Value: aws.String(cfg.Name)},
				},
			},
		},
	}

	out, err := a.client.RunInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	inst := out.Instances[0]
	ip := ""
	if inst.PublicIpAddress != nil {
		ip = *inst.PublicIpAddress
	}

	return &vminterfaces.VirtualMachine{
		ID:	*inst.InstanceId,
		IP:	ip,
	}, nil
}

func (a *AWSProvider) GetVM(id string) (*vminterfaces.VirtualMachine, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	}

	out, err := a.client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return nil, errors.New("instance not found")
	}

	inst := out.Reservations[0].Instances[0]
	ip := ""
	if inst.PublicIpAddress != nil {
		ip = *inst.PublicIpAddress
	}

	return &vminterfaces.VirtualMachine{
		ID:   *inst.InstanceId,
		IP:   ip,
	}, nil
}

func (a *AWSProvider) DeleteVM(id string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	}

	_, err := a.client.TerminateInstances(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (a *AWSProvider) GetVMStatus(instanceID string) (string, error) {
	out, err := a.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return "", err
	}

	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return "", errors.New("instance not found")
	}

	state := out.Reservations[0].Instances[0].State.Name
	return string(state), nil
}