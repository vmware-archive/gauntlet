require 'rails_helper'
require 'json'

describe 'results lifecycle' do

  it 'can create and list results' do
    get '/results'
    expect(JSON.parse(response.body)).to eq('results' => [])

    test_result = {
      pipeline: 'fake-pipe',
      pipecount: 42,
      stage: 'my_stage',
      stagecount: 2,
      jobname: 'my_job',
      gitinfo: 12345,
      pass: true
    }

    post('/results', test_result)
    expect(response.code).to eq('201')
    expect(JSON.parse(response.body)).to eq({ 'resultid' => 1 })

    get '/results'
    resp = JSON.parse(response.body).fetch('results').first
      puts "resp is #{resp}"

    expect(resp['pipeline']).to eq 'fake-pipe'
    expect(resp['pipecount']).to eq 42
    expect(resp['stage']).to eq 'my_stage'
    expect(resp['jobname']).to eq 'my_job'
    expect(resp['pass']).to eq true

    expect(resp['stagecount']).to eq 2
    expect(resp['gitinfo']).to eq '12345'

  end
end