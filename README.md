# Интрукция по запуску
1. Запустить контейнер с RabbitMQ<br>
`docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management`
2. Добавить нового пользователя<br>
`docker exec rabbitmq rabbitmqctl add_user xu senjoxu`
3. Добавить тэги администратора<br>
`docker exec rabbitmq rabbitmqctl set_user_tags xu administrator`
4. Удалить гостя<br>
`docker exec rabbitmq rabbitmqctl delete_user guest`
5. Добавить виртуальный хост<br>
`docker exec rabbitmq rabbitmqctl add_vhost customers`
6. Добавить права на виртуальный хост<br>
`docker exec rabbitmq rabbitmqctl set_permissions -p customers xu ".*" ".*" ".*"`
7. Добавить exchange<br>
`docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=fanout -u xu -p senjoxu durable=true`
8. Добавить права на exchange<br>
`docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers xu customer_events ".*" ".*"`
9. Добавить exchange с колбэками<br>
`docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_callbacks type=direct -u xu -p senjoxu durable=true`
10. Добавить права на exchange с колбэками<br>
`docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers xu customer_callbacks ".*" ".*"`
