Gauntlet::Application.routes.draw do
  get 'results' => 'results#index'
  post 'results' => 'results#create'
end
