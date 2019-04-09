from locust import HttpLocust, TaskSet

def reportIncident(l):
    l.client.post("/rest/log", {"Reporter":"0xFE00BB37A56282d33680542Ae1CD6763660b5555","Message":"automatic reporting", "Location":"NYC-Datacenter:VM12345"})

def getLogs(l):
    l.client.get("/logs")

def index(l):
    l.client.get("/")

class UserBehavior(TaskSet):
    tasks = {index: 1, reportIncident: 2, getLogs: 1}

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000