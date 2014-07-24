require 'rails_helper'

describe ResultsController do
  describe 'POST #create result' do
    let(:test_result) do
      {
        pipeline: 'fake-pipe',
        pipecount: 42,
        stage: 'my_stage',
        stagecount: 2,
        jobname: 'my_job',
        gitinfo: 12345,
        pass: true
      }
    end

    it 'creates a result' do
      expect { post :create, test_result }.to change { Result.count }.by(1)
    end

    it 'returns the result ID' do
      first_response = post :create, test_result
      expect(JSON.parse(first_response.body)).to eq( {'resultid' => 1} )

      second_response = post :create, test_result
      expect(JSON.parse(second_response.body)).to eq( {'resultid' => 2} )
    end
  end
end