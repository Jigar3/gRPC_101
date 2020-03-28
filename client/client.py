import grpc

import sys
sys.path.insert(1, '/home/jigar/grpc')


import githubpb.github_pb2 as a 
import githubpb.github_pb2_grpc as b

channel = grpc.insecure_channel("localhost:50051")

stub = b.GithubServiceStub(channel)

# number = a.SumRequest(first_number=4, second_number=4)
user = a.FollowerRequest(github_username="jigar3")

res = stub.GetFollowers(user)

print(res.follower_list)
