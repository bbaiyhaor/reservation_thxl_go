import 'weui';
import { MediaBox, MediaBoxBody, MediaBoxDescription, MediaBoxHeader, MediaBoxTitle, Panel, PanelBody, PanelHeader } from '#react-weui';
import React, { PropTypes } from 'react';
import PageBottom from '#coms/PageBottom';
import ReactDOM from 'react-dom';
import StudentEntryIcon from '#imgs/mobile/student.png';
import TeacherEntryIcon from '#imgs/mobile/teacher.png';

function EntryApp({ studentEntryIcon, teacherEntryIcon }) {
  return (
    <Panel>
      <PanelHeader style={{ fontSize: '18px' }}>
        学生心理发展指导中心预约系统
      </PanelHeader>
      <PanelBody>
        <MediaBox type="appmsg" href="/m/student">
          <MediaBoxHeader>
            <img src={studentEntryIcon} alt="Student Entry Icon" />
          </MediaBoxHeader>
          <MediaBoxBody>
            <MediaBoxTitle>我是学生</MediaBoxTitle>
            <MediaBoxDescription>点击进入</MediaBoxDescription>
          </MediaBoxBody>
        </MediaBox>
        <MediaBox type="appmsg" href="/m/teacher">
          <MediaBoxHeader>
            <img src={teacherEntryIcon} alt="Teacher Entry Icon" />
          </MediaBoxHeader>
          <MediaBoxBody>
            <MediaBoxTitle>我是咨询师</MediaBoxTitle>
            <MediaBoxDescription>点击进入</MediaBoxDescription>
          </MediaBoxBody>
        </MediaBox>
      </PanelBody>
      <PageBottom
        styles={{ color: '#999999', textAlign: 'center' }}
        contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
      />
    </Panel>
  );
}

EntryApp.propTypes = {
  studentEntryIcon: PropTypes.string.isRequired,
  teacherEntryIcon: PropTypes.string.isRequired,
};

ReactDOM.render(
  <EntryApp
    studentEntryIcon={StudentEntryIcon}
    teacherEntryIcon={TeacherEntryIcon}
  />,
  document.getElementById('content'),
);
