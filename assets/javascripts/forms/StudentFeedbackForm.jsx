import 'weui';
import { Button, ButtonArea, Cell, CellBody, CellFooter, Cells, CellsTitle, Form, FormCell, Radio } from 'react-weui';
import React, { PropTypes } from 'react';

export default class StudentFeedbackForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      feedbackRadio1: 0,
      feedbackRadio2: 0,
      feedbackRadio3: 0,
      feedbackRadio4: 0,
      feedbackRadio5: 0,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { scores } = nextProps;
    if (scores && scores.length === 5) {
      this.setState({
        feedbackRadio1: nextProps.scores[0],
        feedbackRadio2: nextProps.scores[1],
        feedbackRadio3: nextProps.scores[2],
        feedbackRadio4: nextProps.scores[3],
        feedbackRadio5: nextProps.scores[4],
      });
    } else {
      this.setState({
        feedbackRadio1: 4,
        feedbackRadio2: 4,
        feedbackRadio3: 4,
        feedbackRadio4: 4,
        feedbackRadio5: 4,
      });
    }
  }

  handleChange(e) {
    this.setState({ [e.target.name]: Number(e.target.value) });
  }

  handleSubmit() {
    this.props.handleSubmit(this.state.feedbackRadio1, this.state.feedbackRadio2, this.state.feedbackRadio3, this.state.feedbackRadio4, this.state.feedbackRadio5);
  }

  render() {
    return (
      <div>
        {this.props.reservation &&
          <CellsTitle>
            正在反馈：{this.props.reservation.start_time}-{this.props.reservation.end_time.slice(-5)} {this.props.reservation.teacher_fullname}
          </CellsTitle>
        }
        <Form radio>
          <Cells>
            <Cell>
              <CellBody>
                1、你是否得到了你所希望的咨询？
              </CellBody>
            </Cell>
            <FormCell radio>
              <CellBody>肯定是的</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio1"
                  value="4"
                  checked={this.state.feedbackRadio1 === 4}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>基本上是的</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio1"
                  value="3"
                  checked={this.state.feedbackRadio1 === 3}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>没有</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio1"
                  value="2"
                  checked={this.state.feedbackRadio1 === 2}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>肯定没有</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio1"
                  value="1"
                  checked={this.state.feedbackRadio1 === 1}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
          </Cells>
          <Cells>
            <Cell>
              <CellBody>
                2、咨询在多大程度上满足了你的需要？
              </CellBody>
            </Cell>
            <FormCell radio>
              <CellBody>几乎全部需要得到满足</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio2"
                  value="4"
                  checked={this.state.feedbackRadio2 === 4}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>大部分需要得到满足</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio2"
                  value="3"
                  checked={this.state.feedbackRadio2 === 3}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>仅一小部分需要得到满足</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio2"
                  value="2"
                  checked={this.state.feedbackRadio2 === 2}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>需要丝毫没有得到满足</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio2"
                  value="1"
                  checked={this.state.feedbackRadio2 === 1}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
          </Cells>
          <Cells>
            <Cell>
              <CellBody>
                3、如果一个朋友需要咨询，你会向他或她推荐这位咨询师吗？
              </CellBody>
            </Cell>
            <FormCell radio>
              <CellBody>肯定会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio3"
                  value="4"
                  checked={this.state.feedbackRadio3 === 4}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio3"
                  value="3"
                  checked={this.state.feedbackRadio3 === 3}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>不会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio3"
                  value="2"
                  checked={this.state.feedbackRadio3 === 2}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>肯定不会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio3"
                  value="1"
                  checked={this.state.feedbackRadio3 === 1}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
          </Cells>
          <Cells>
            <Cell>
              <CellBody>
                4、总体来讲，你对你接受的咨询有多满意？
              </CellBody>
            </Cell>
            <FormCell radio>
              <CellBody>非常满意</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio4"
                  value="4"
                  checked={this.state.feedbackRadio4 === 4}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>大部分满意</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio4"
                  value="3"
                  checked={this.state.feedbackRadio4 === 3}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>无所谓，或不太满意</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio4"
                  value="2"
                  checked={this.state.feedbackRadio4 === 2}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>非常不满意</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio4"
                  value="1"
                  checked={this.state.feedbackRadio4 === 1}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
          </Cells>
          <Cells>
            <Cell>
              <CellBody>
                5、如果你将再次寻求咨询，你会回来找这位咨询师吗？
              </CellBody>
            </Cell>
            <FormCell radio>
              <CellBody>肯定会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio5"
                  value="4"
                  checked={this.state.feedbackRadio5 === 4}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio5"
                  value="3"
                  checked={this.state.feedbackRadio5 === 3}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>不会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio5"
                  value="2"
                  checked={this.state.feedbackRadio5 === 2}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
            <FormCell radio>
              <CellBody>肯定不会</CellBody>
              <CellFooter>
                <Radio
                  name="feedbackRadio5"
                  value="1"
                  checked={this.state.feedbackRadio5 === 1}
                  onChange={this.handleChange}
                />
              </CellFooter>
            </FormCell>
          </Cells>
        </Form>
        <ButtonArea direction="horizontal">
          <Button onClick={this.handleSubmit}>提交反馈</Button>
          <Button type="default" onClick={this.props.handleCancel}>返回</Button>
        </ButtonArea>
        <div style={{ height: '10px' }} />
      </div>
    );
  }
}

StudentFeedbackForm.propTypes = {
  reservation: PropTypes.object.isRequired,
  scores: PropTypes.arrayOf(PropTypes.number).isRequired,
  handleSubmit: PropTypes.func.isRequired,
  handleCancel: PropTypes.func.isRequired,
};
