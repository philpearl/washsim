# Washsim

Simulating either cars in a car-wash or CPU-bound tasks on a server, depending
on your point of view. Heavily informed by https://simpy.readthedocs.io/en/latest/examples/carwash.html,
Discussed at https://medium.com/@philpearl/carwash-simulations-cpu-usage-12a19fe25bf#.qqavo2m62.

The obvious limitation to this simulation is that it runs in real-time. The next
thing to attempt would be to fake time - incrementing it instantly to the next 
time that an event would occur.

## License
MIT licence