import random 

def gen_random_es(maxsize, mem_prob):
    ms = "values :=map[uint64]bool{\n"
    for i in range(maxsize):
        r = random.random()
        if r <= mem_prob:
            ms += f"{i} : true,\n"
        else:
            ms += f"{i} : false,\n"
    ms += "}\n"
    return ms

print(gen_random_es(32,0.25))