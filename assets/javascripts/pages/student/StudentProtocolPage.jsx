import 'weui';
import { Article } from '#react-weui';
import PageBottom from '#coms/PageBottom';
import React from 'react';

export default class StudentProtocolPage extends React.Component {
  componentDidMount() {
    window.scroll(0, 0);
  }

  render() {
    const styles = {
      h1: {
        textAlign: 'center',
        fontWeight: 'bold',
        fontSize: '22px',
      },
      h2: {
        fontWeight: 'bold',
        fontSize: '18px',
      },
    };

    return (
      <div>
        <Article>
          <h1 style={{ ...styles.h1 }}>心理发展指导中心咨询协议</h1>
          <section>
            <h2 className="title" style={{ ...styles.h2 }}>1．咨询对象</h2>
            <section>
              <p>　　清华大学全日制在校学生。首次来访请携带学生证。</p>
            </section>
            <h2 className="title" style={{ ...styles.h2 }}>2．咨询内容</h2>
            <section>
              <p>　　心理发展指导中心主要提供心理发展、心理适应及心理障碍等方面的咨询。</p>
            </section>
            <h2 className="title" style={{ ...styles.h2 }}>3. 咨询原则</h2>
            <section>
              <p>　　本中心尊重来访者的个人隐私，任何个人信息都不会被泄露，除非有你本人的授权。
                <b>但在你的信息暗示你将危及自己或他人及社会的安全时，咨询师有权和有关方面联系。</b>
              </p>
            </section>
            <section>
              <p>　　另外，为提高咨询水平，保障来访者接受到高质量的咨询服务，有时需在咨询过程中录音或者录像，
                用于个案督导、个案研讨或教学研究。所有音像资料只限于专业人员，决不会有非专业人员涉及。
                如有录音、录像的要求，咨询师会事前征得你的同意。
              </p>
            </section>
            <h2 className="title" style={{ ...styles.h2 }}>4. 来访者的职责</h2>
            <section>
              <p>　　心理咨询的过程，既有收益又有风险。收益是通过咨询排解心理困扰、开发心理潜能、促进个人成长。
                咨询的风险包括记起不愉快的往事，激发起很强烈的情绪等。在咨询过程中保持积极、开放和诚实的态度非常重要，
                你主要的职责就是向着你和咨询师共同制定的目标而努力。
                <b>如果在咨询过程中，受到侵犯和伤害，你有权随时提出中断咨询，更换咨询师。</b>
                中心工作人员将及时安排新的咨询师或者帮助你转介到其它机构。
              </p>
            </section>
            <h2 className="title" style={{ ...styles.h2 }}>5. 对来访者的要求</h2>
            <section>
              <p>　　学生可以通过网络了解心理发展指导中心及咨询师的基本情况，并通过网络和电话预约咨询。
                如有预约，应按预约时间准时到达。不能如期来访，应提前打电话说明（62782007）；
                如无预约，应听从咨询中心的安排进行咨询。
              </p>
            </section>
            <section>
              <p>　　一般每次咨询不超过50分钟。一次咨询不能完成者，咨询师将再约时间咨询。</p>
            </section>
            <section>
              <p>　　咨询过程需在咨询室内完成，来访者不得要求在咨询室以外进行心理咨询。</p>
            </section>
            <section>
              <p>　　在等待咨询的过程中保持室内安静；爱护咨询室的公共设施。</p>
            </section>
          </section>
        </Article>
        <PageBottom
          styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
          contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
          height="55px"
        />
      </div>
    );
  }
}
