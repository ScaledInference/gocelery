from celery import Celery


# app = Celery('tasks',
#     broker='ampq://localhost:6379',
#     backend='redis://localhost:6379'
# )
#
app = Celery('example', backend='amqp://', broker='amqp://')

app.conf.update(
    CELERY_TASK_SERIALIZER='json',
    CELERY_ACCEPT_CONTENT=['json'],  # Ignore other content
    CELERY_RESULT_SERIALIZER='json',
    CELERY_ENABLE_UTC=True,
)


@app.task
def add(x, y):
    print(f'ADDING{x}, {y}')
    return x + y


@app.task
def add_reflect(x, y):
    return x + y
