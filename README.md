api/user/register/ 
api/user/login post отправляем логин, пароль 


api/portfolio/ post создаём портфель
            get получаем список портфелей

api/portfolio/{id_portfolio} get получаем состав портфеля, с текущим курсом

api/portfolio/{id_portfolio} post добавляем или убавляем монету в портфель 

api/portfolio/{id_portfolio}/{id_coin}/ get история действий с монетой

api/portfolio/{id_portfolio}/{id_coin}/{id_record} delete удаление действия 

